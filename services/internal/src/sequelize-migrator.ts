import fs from "fs";
import path from "path";
import { pathToFileURL } from "url";
import { Sequelize } from "sequelize";
import type { MigrationModule } from "./sequelize-migration";

interface MigrationEntry {
    canonicalName: string;
    migrationPath: string;
}

function getDatabaseUrl(): string {
    return process.env.DATABASE_URL || "postgres://fm:admin@127.0.0.1:5432/fm";
}

function getDialectOptions(): Record<string, unknown> {
    if (process.env.PGSSL !== "1") {
        return {};
    }
    return { ssl: { require: true, rejectUnauthorized: false } };
}

function getMigrationsDir(): string {
    return path.join(__dirname, "..", "sequelize", "migrations");
}

function toCanonicalName(fileName: string): string {
    return fileName.endsWith(".ts") ? `${fileName.slice(0, -3)}.js` : fileName;
}

function readMigrationEntries(): MigrationEntry[] {
    const migrationsDir = getMigrationsDir();
    const files = fs
        .readdirSync(migrationsDir)
        .filter((file) => file.endsWith(".ts") || file.endsWith(".js"))
        .sort();

    const byCanonical = new Map<string, MigrationEntry>();
    for (const file of files) {
        const canonicalName = toCanonicalName(file);
        const migrationPath = path.join(migrationsDir, file);
        const existing = byCanonical.get(canonicalName);
        if (!existing || file.endsWith(".ts")) {
            byCanonical.set(canonicalName, { canonicalName, migrationPath });
        }
    }
    return Array.from(byCanonical.values()).sort((a, b) => a.canonicalName.localeCompare(b.canonicalName));
}

async function ensureSequelizeMeta(sequelize: Sequelize): Promise<void> {
    await sequelize.query(`
        CREATE TABLE IF NOT EXISTS "SequelizeMeta" (
            name VARCHAR(255) PRIMARY KEY
        );
    `);
}

async function loadMigrationModule(migrationPath: string): Promise<MigrationModule> {
    const loaded = await import(pathToFileURL(migrationPath).href);
    const migration = (loaded.default ?? loaded) as Partial<MigrationModule>;
    if (!migration.up || !migration.down) {
        throw new Error(`Invalid migration module: ${migrationPath}`);
    }
    return migration as MigrationModule;
}

export async function runMigrate(): Promise<void> {
    const sequelize = new Sequelize(getDatabaseUrl(), {
        dialect: "postgres",
        logging: false,
        dialectOptions: getDialectOptions(),
    });

    try {
        await ensureSequelizeMeta(sequelize);
        const [appliedRowsRaw] = await sequelize.query('SELECT name FROM "SequelizeMeta" ORDER BY name');
        const appliedRows = appliedRowsRaw as Array<{ name: string }>;
        const applied = new Set(appliedRows.map((row) => row.name));
        const queryInterface = sequelize.getQueryInterface();

        for (const entry of readMigrationEntries()) {
            if (applied.has(entry.canonicalName)) {
                continue;
            }
            const migration = await loadMigrationModule(entry.migrationPath);
            await migration.up(queryInterface, await import("sequelize"));
            await sequelize.query('INSERT INTO "SequelizeMeta" (name) VALUES ($1)', {
                bind: [entry.canonicalName],
            });
            console.log(`[migrate] applied ${entry.canonicalName}`);
        }
    } finally {
        await sequelize.close();
    }
}

export async function runMigrateUndo(): Promise<void> {
    const sequelize = new Sequelize(getDatabaseUrl(), {
        dialect: "postgres",
        logging: false,
        dialectOptions: getDialectOptions(),
    });

    try {
        await ensureSequelizeMeta(sequelize);
        const [rowsRaw] = await sequelize.query('SELECT name FROM "SequelizeMeta" ORDER BY name DESC LIMIT 1');
        const rows = rowsRaw as Array<{ name: string }>;
        const lastApplied = rows[0]?.name;
        if (!lastApplied) {
            console.log("[migrate:undo] no applied migrations");
            return;
        }

        const migration = readMigrationEntries().find((entry) => entry.canonicalName === lastApplied);
        if (!migration) {
            throw new Error(`Migration file not found for ${lastApplied}`);
        }

        const queryInterface = sequelize.getQueryInterface();
        const migrationModule = await loadMigrationModule(migration.migrationPath);
        await migrationModule.down(queryInterface, await import("sequelize"));
        await sequelize.query('DELETE FROM "SequelizeMeta" WHERE name = $1', {
            bind: [lastApplied],
        });
        console.log(`[migrate:undo] reverted ${lastApplied}`);
    } finally {
        await sequelize.close();
    }
}
