"use strict";

const path = require("path");
const fs = require("fs/promises");

class WzService {
    constructor(appConfiguration) {
        this.wzRoot = appConfiguration?.resources?.wz_root;
        this.itemMetaById = new Map();
    }

    _parseEnhanceChance(infoBlock) {
        const tuc = infoBlock.match(/<(?:int|short) name="tuc" value="(-?\d+)"/);
        if (!tuc) {
            return 0;
        }
        return Math.max(0, Number(tuc[1]));
    }

    _upsertItem(id, infoBlock) {
        const itemId = Number(id);
        if (!Number.isInteger(itemId) || itemId <= 0) {
            return;
        }
        this.itemMetaById.set(itemId, {
            enhanceChance: this._parseEnhanceChance(infoBlock),
        });
    }

    async _loadFile(filePath) {
        const xml = await fs.readFile(filePath, "utf8");

        const singleRoot = xml.match(/<imgdir name="0?(\d{7})\.img">/);
        if (singleRoot) {
            const infoBlock = xml.match(/<imgdir name="info">([\s\S]*?)<\/imgdir>/);
            this._upsertItem(singleRoot[1], infoBlock ? infoBlock[1] : "");
            return;
        }

        const itemRegex = /<imgdir name="(\d{7})">[\s\S]*?<imgdir name="info">([\s\S]*?)<\/imgdir>/g;
        for (const m of xml.matchAll(itemRegex)) {
            this._upsertItem(m[1], m[2]);
        }
    }

    async _collectXmlFiles(rootDir) {
        const files = [];
        const queue = [rootDir];
        while (queue.length > 0) {
            const dir = queue.pop();
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

    async preload() {
        const targets = await this._collectXmlFiles(this.wzRoot);
        const batchSize = 32;
        for (let i = 0; i < targets.length; i += batchSize) {
            const batch = targets.slice(i, i + batchSize);
            await Promise.all(batch.map((p) => this._loadFile(p)));
        }
    }

    async getEnhanceChance(itemId) {
        const id = Number(itemId);
        if (!Number.isInteger(id) || id <= 0) {
            return 0;
        }
        return this.itemMetaById.get(id)?.enhanceChance ?? 0;
    }
}

module.exports = { WzService };
