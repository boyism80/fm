import path from "path";
import fs from "fs/promises";
import type { AppConfiguration } from "../config/app-configuration";

const MAKE_CHAR_INFO_GENDER: Record<string, number> = { CharMale: 0, CharFemale: 1 };

export class WzService {
    private readonly wzRoot: string;
    private readonly itemMetaById: Map<number, { enhanceChance: number }>;
    // gender -> [face, hair, top, bottom, shoes, weapon] allowed ids
    private readonly makeCharInfo: Map<number, Set<number>[]>;
    private readonly forbiddenNames: string[];

    constructor(appConfiguration: AppConfiguration) {
        this.wzRoot = appConfiguration?.resources?.wz_root;
        this.itemMetaById = new Map();
        this.makeCharInfo = new Map();
        this.forbiddenNames = [];
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
        await this.loadMakeCharInfo(path.join(this.wzRoot, "Etc.wz", "MakeCharInfo.img.xml"));
        await this.loadForbiddenNames(path.join(this.wzRoot, "Etc.wz", "ForbiddenName.img.xml"));
    }

    private async loadMakeCharInfo(filePath: string): Promise<void> {
        const xml = await fs.readFile(filePath, "utf8");
        const stack: string[] = [];
        for (const m of xml.matchAll(/<imgdir name="([^"]+)">|<\/imgdir>|<int name="\d+" value="(-?\d+)"\/>/g)) {
            if (m[1] !== undefined) {
                stack.push(m[1]);
                continue;
            }
            if (m[2] === undefined) {
                stack.pop();
                continue;
            }
            // MakeCharInfo.img / CharMale / {part}
            if (stack.length !== 3) {
                continue;
            }
            const gender = MAKE_CHAR_INFO_GENDER[stack[1] ?? ""];
            const part = Number(stack[2]);
            if (gender === undefined || Number.isInteger(part) === false || part < 0 || part > 5) {
                continue;
            }
            let parts = this.makeCharInfo.get(gender);
            if (parts == null) {
                parts = [new Set(), new Set(), new Set(), new Set(), new Set(), new Set()];
                this.makeCharInfo.set(gender, parts);
            }
            parts[part]?.add(Number(m[2]));
        }
    }

    private async loadForbiddenNames(filePath: string): Promise<void> {
        const xml = await fs.readFile(filePath, "utf8");
        for (const m of xml.matchAll(/<string name="\d+" value="([^"]+)"\/>/g)) {
            this.forbiddenNames.push((m[1] ?? "").toLowerCase());
        }
    }

    /** face, hair, top, bottom, shoes and weapon must all be in the gender's MakeCharInfo lists. */
    canMakeCharacter(gender: number, looks: [number, number, number, number, number, number]): boolean {
        const parts = this.makeCharInfo.get(gender);
        if (parts == null) {
            return false;
        }
        return looks.every((id, idx) => parts[idx]?.has(id) === true);
    }

    isForbiddenName(name: string): boolean {
        const lower = name.toLowerCase();
        return this.forbiddenNames.some((word) => word.length > 0 && lower.includes(word));
    }

    async getEnhanceChance(itemId: number): Promise<number> {
        if (!Number.isInteger(itemId) || itemId <= 0) {
            return 0;
        }
        return this.itemMetaById.get(itemId)?.enhanceChance ?? 0;
    }
}
