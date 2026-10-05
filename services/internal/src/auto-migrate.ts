import { applyPendingMigrations, connectMigrationDatabase } from "./sequelize-migrator";
import type { InternalConfig, PgEndpoint } from "./types/internal-config";

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

export async function autoMigrateAllIfEnabled(internalConfig: InternalConfig): Promise<void> {
    const enabled = internalConfig.sequelize?.auto_migrate_on_startup ?? true;
    if (!enabled) {
        return;
    }
    for (const target of collectPgEndpoints(internalConfig)) {
        const sequelize = connectMigrationDatabase(endpointToUrl(target.endpoint), Boolean(target.endpoint.ssl));
        try {
            await applyPendingMigrations(sequelize, `[auto-migrate] ${target.name}:`);
        } finally {
            await sequelize.close();
        }
    }
}
