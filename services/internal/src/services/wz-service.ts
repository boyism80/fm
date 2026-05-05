import path from "path";
import fs from "fs/promises";
import type { AppConfiguration } from "../config/app-configuration";

export class WzService {
    private readonly wzRoot: string;
    private readonly itemMetaById: Map<number, { enhanceChance: number }>;

    constructor(appConfiguration: AppConfiguration) {
        this.wzRoot = appConfiguration?.resources?.wz_root;
        this.itemMetaById = new Map();
    }

    private parseEnhanceChance(infoBlock: string): number {
        const tuc = infoBlock.match(/<(?:int|short) name="tuc" value="(-?\d+)"/);
        if (!tuc) {
            return 0;
        }
        return Math.max(0, Number(tuc[1]));
    }

    private upsertItem(id: string, infoBlock: string): void {
        const itemId = Number(id);
        if (!Number.isInteger(itemId) || itemId <= 0) {
            return;
        }
        this.itemMetaById.set(itemId, { enhanceChance: this.parseEnhanceChance(infoBlock) });
    }

    private async loadFile(filePath: string): Promise<void> {
        const xml = await fs.readFile(filePath, "utf8");
        const singleRoot = xml.match(/<imgdir name="0?(\d{7})\.img">/);
        if (singleRoot) {
            const infoBlock = xml.match(/<imgdir name="info">([\s\S]*?)<\/imgdir>/);
            const rootId = singleRoot[1];
            this.upsertItem(rootId ?? "", infoBlock ? (infoBlock[1] ?? "") : "");
            return;
        }
        const itemRegex = /<imgdir name="(\d{7})">[\s\S]*?<imgdir name="info">([\s\S]*?)<\/imgdir>/g;
        for (const m of xml.matchAll(itemRegex)) {
            this.upsertItem(m[1] ?? "", m[2] ?? "");
        }
    }

    private async collectXmlFiles(rootDir: string): Promise<string[]> {
        const files: string[] = [];
        const queue = [rootDir];
        while (queue.length > 0) {
            const dir = queue.pop();
            if (!dir) {
                continue;
            }
            const entries = await fs.readdir(dir, { withFileTypes: true });
            for (const entry of entries) {
                const fullPath = path.join(dir, entry.name);
                if (entry.isDirectory()) {
                    queue.push(fullPath);
                } else if (entry.isFile() && entry.name.endsWith(".img.xml")) {
                    files.push(fullPath);
                }
            }
        }
        return files;
    }

    async preload(): Promise<void> {
        const targets = await this.collectXmlFiles(this.wzRoot);
        const batchSize = 32;
        for (let i = 0; i < targets.length; i += batchSize) {
            const batch = targets.slice(i, i + batchSize);
            await Promise.all(batch.map((p) => this.loadFile(p)));
        }
    }

    async getEnhanceChance(itemId: number): Promise<number> {
        if (!Number.isInteger(itemId) || itemId <= 0) {
            return 0;
        }
        return this.itemMetaById.get(itemId)?.enhanceChance ?? 0;
    }
}
