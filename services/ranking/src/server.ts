import fs from "node:fs";
import os from "node:os";
import { fileURLToPath } from "node:url";
import yaml from "js-yaml";
import { autoMigrateAllIfEnabled } from "@fm/shared/auto-migrate";
import { loadConfig } from "@fm/shared/config";
import { AppConfiguration } from "@fm/shared/config/app-configuration";
import { InternalContext } from "@fm/shared/context/internal-context";
import { DailyLogFile } from "@fm/shared/log-file";
import { RankingRepository } from "@fm/shared/repos/ranking-repository";
import { RankingScheduler, type RankingSchedulerOptions } from "./ranking-scheduler";

type RawRankingConfig = Partial<{
    daily_hour: number;
    range_multiplier: number;
    workers: number;
    world_groups: Array<{ id: number; worlds: number[] }>;
}>;

function loadOptions(configPath: string): RankingSchedulerOptions {
    const doc = (yaml.load(fs.readFileSync(configPath, "utf8")) ?? {}) as { ranking?: RawRankingConfig };
    const raw = doc.ranking ?? {};
    const worldGroups = new Map<number, number>();
    for (const group of raw.world_groups ?? []) {
        for (const world of group.worlds) {
            worldGroups.set(world, group.id);
        }
    }
    const options = {
        dailyHour: raw.daily_hour ?? 1,
        rangeMultiplier: raw.range_multiplier ?? 4,
        workers: raw.workers ?? os.cpus().length,
        worldGroups,
    };
    if (!Number.isInteger(options.dailyHour) || options.dailyHour < 0 || options.dailyHour > 23) {
        throw new Error(`ranking.daily_hour must be an hour from 0 to 23: ${options.dailyHour}`);
    }
    if (!Number.isInteger(options.rangeMultiplier) || options.rangeMultiplier < 1) {
        throw new Error(`ranking.range_multiplier must be a positive integer: ${options.rangeMultiplier}`);
    }
    if (!Number.isInteger(options.workers) || options.workers < 1) {
        throw new Error(`ranking.workers must be a positive integer: ${options.workers}`);
    }
    return options;
}

async function main() {
    new DailyLogFile("ranking").captureConsole();
    const configPath = process.env.FM_RANKING_CONFIG ?? fileURLToPath(new URL("../config.yaml", import.meta.url));
    const config = loadConfig(configPath);
    const options = loadOptions(configPath);
    await autoMigrateAllIfEnabled(config);

    const internalContext = new InternalContext(new AppConfiguration(config));
    const scheduler = new RankingScheduler(internalContext, new RankingRepository(internalContext), options);
    await scheduler.start();
    console.log(
        `fm ranking scheduler started (daily refresh ${options.dailyHour}:00 KST, workers ${options.workers}, range multiplier ${options.rangeMultiplier})`,
    );

    async function shutdown() {
        setTimeout(() => {
            internalContext.close().catch(() => {});
            process.exit(1);
        }, 10_000).unref();
        await scheduler.stop();
        await internalContext.close().catch(() => {});
        process.exit(0);
    }

    process.on("SIGINT", shutdown);
    process.on("SIGTERM", shutdown);
}

main().catch((err: unknown) => {
    console.error("[startup] fatal:", err);
    process.exit(1);
});
