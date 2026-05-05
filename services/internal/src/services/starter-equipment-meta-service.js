"use strict";

const path = require("path");
const fs = require("fs/promises");

class StarterEquipmentMetaService {
    constructor(appConfiguration) {
        const wzRoot = appConfiguration?.resources?.wz_root;
        this.characterWzRoot = path.join(wzRoot, "Character.wz");
        this.byItemId = new Map();
    }

    async _loadFile(filePath) {
        const xml = await fs.readFile(filePath, "utf8");

        const singleRoot = xml.match(/<imgdir name="0?(\d{7})\.img">/);
        if (singleRoot) {
            const id = Number(singleRoot[1]);
            const infoBlock = xml.match(/<imgdir name="info">([\s\S]*?)<\/imgdir>/);
            const info = infoBlock ? infoBlock[1] : "";
            const tuc = info.match(/<(?:int|short) name="tuc" value="(-?\d+)"/);
            this.byItemId.set(id, tuc ? Math.max(0, Number(tuc[1])) : 0);
        } else {
            const itemRegex = /<imgdir name="(\d{7})">[\s\S]*?<imgdir name="info">([\s\S]*?)<\/imgdir>/g;
            for (const m of xml.matchAll(itemRegex)) {
                const id = Number(m[1]);
                const infoBlock = m[2];
                const tuc = infoBlock.match(/<(?:int|short) name="tuc" value="(-?\d+)"/);
                this.byItemId.set(id, tuc ? Math.max(0, Number(tuc[1])) : 0);
            }
        }
    }

    async preload() {
        const folders = [
            "Cap",
            "Accessory",
            "Coat",
            "Longcoat",
            "Pants",
            "Shoes",
            "Glove",
            "Shield",
            "Cape",
            "Ring",
            "Weapon",
        ];
        for (const folder of folders) {
            const dirPath = path.join(this.characterWzRoot, folder);
            let files = [];
            try {
                files = await fs.readdir(dirPath, { withFileTypes: true });
            } catch (err) {
                console.error(`[starter-equipment-meta] preload readdir failed: ${dirPath}`, err);
                throw err;
            }
            const targets = files
                .filter((d) => d.isFile() && d.name.endsWith(".img.xml"))
                .map((d) => path.join(dirPath, d.name));
            await Promise.all(targets.map((p) => this._loadFile(p)));
        }
    }

    async getEnhanceChance(itemId) {
        const id = Number(itemId);
        if (!Number.isInteger(id) || id <= 0) {
            return 0;
        }
        if (this.byItemId.size === 0) {
            console.error(
                `[starter-equipment-meta] cache empty in getEnhanceChance: itemId=${id}. preload may have failed or was not run`
            );
            return 0;
        }
        return this.byItemId.get(id) ?? 0;
    }
}

module.exports = { StarterEquipmentMetaService };
