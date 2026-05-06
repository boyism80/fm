import fs from "fs";
import path from "path";
import { createRequire } from "module";
import { fileURLToPath } from "url";
import { DataTypes, Sequelize } from "sequelize";
import type { InternalConfig, PgEndpoint } from "./types/internal-config";

type MigrationQueryInterface = ReturnType<Sequelize["getQueryInterface"]>;
const require = createRequire(import.meta.url);
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

interface TargetEndpoint {
    name: string;
    endpoint: PgEndpoint;
}

function endpointToUrl(endpoint: PgEndpoint): string {
    const user = encodeURIComponent(endpoint.username ?? "");
    const pass = encodeURIComponent(endpoint.password ?? "");
    const auth = pass ? `${user}:${pass}` : user;
    return `postgres://${auth}@${endpoint.host}:${endpoint.port}/${endpoint.database}`;
}

function collectPgEndpoints(config: InternalConfig): TargetEndpoint[] {
    const targets: TargetEndpoint[] = [];
    if (config.postgresql.unified) {
        targets.push({ name: "unified", endpoint: config.postgresql.unified });
    }
    for (const [worldId, world] of Object.entries(config.postgresql.worlds ?? {})) {
        targets.push({ name: `world-${worldId}-global`, endpoint: world.global });
        for (let i = 0; i < (world.data ?? []).length; i += 1) {
            const shard = world.data[i];
            if (shard) {
                targets.push({ name: `world-${worldId}-data-${i}`, endpoint: shard });
            }
        }
    }
    return targets;
}

async function ensureSequelizeMeta(queryInterface: MigrationQueryInterface): Promise<void> {
    const [rowsRaw] = await queryInterface.sequelize.query(
        "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'SequelizeMeta') AS exists"
    );
    const rows = rowsRaw as Array<{ exists?: boolean }>;
    if (!rows[0]?.exists) {
        await queryInterface.createTable("SequelizeMeta", {
            name: {
                type: DataTypes.STRING,
                allowNull: false,
                primaryKey: true,
            },
        });
    }
}

async function runMigrationsOnEndpoint(baseDir: string, target: TargetEndpoint): Promise<void> {
    const migrationsDir = path.join(baseDir, "sequelize", "migrations");
    const files = fs.readdirSync(migrationsDir).filter((file) => file.endsWith(".ts") || file.endsWith(".js")).sort();
    const sequelize = new Sequelize(endpointToUrl(target.endpoint), {
        dialect: "postgres",
        logging: false,
        dialectOptions: target.endpoint.ssl ? { ssl: { require: true, rejectUnauthorized: false } } : {},
    });
    try {
        const queryInterface = sequelize.getQueryInterface();
        await ensureSequelizeMeta(queryInterface);
        const [appliedRowsRaw] = await sequelize.query("SELECT name FROM \"SequelizeMeta\" ORDER BY name");
        const appliedRows = appliedRowsRaw as Array<{ name: string }>;
        const applied = new Set(appliedRows.map((row) => row.name));
        const canonicalNames = new Set<string>();
        for (const file of files) {
            const canonicalName = file.endsWith(".ts") ? `${file.slice(0, -3)}.js` : file;
            if (canonicalNames.has(canonicalName)) {
                continue;
            }
            canonicalNames.add(canonicalName);
            if (applied.has(canonicalName)) {
                continue;
            }
            const migrationPath = path.join(migrationsDir, file);
            delete require.cache[require.resolve(migrationPath)];
            const migrationModule = require(migrationPath) as {
                up?: (qi: MigrationQueryInterface, sq: typeof Sequelize) => Promise<void>;
                default?: { up?: (qi: MigrationQueryInterface, sq: typeof Sequelize) => Promise<void> };
            };
            const migration = migrationModule.default ?? migrationModule;
            if (typeof migration.up !== "function") {
                throw new Error(`Invalid migration module (missing up): ${migrationPath}`);
            }
            await migration.up(queryInterface, Sequelize);
            await sequelize.query("INSERT INTO \"SequelizeMeta\" (name) VALUES ($1)", { bind: [canonicalName] });
            console.log(`[auto-migrate] ${target.name}: applied ${canonicalName}`);
        }
    } finally {
        await sequelize.close();
    }
}

export async function autoMigrateAllIfEnabled(internalConfig: InternalConfig): Promise<void> {
    const enabled = internalConfig.sequelize?.auto_migrate_on_startup ?? false;
    if (!enabled) {
        return;
    }
    const baseDir = path.join(__dirname, "..");
    const targets = collectPgEndpoints(internalConfig);
    for (const target of targets) {
        await runMigrationsOnEndpoint(baseDir, target);
    }
}
