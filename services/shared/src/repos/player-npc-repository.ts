import type { InternalContext } from "../context/internal-context";

export type PlayerNpcEquipRow = {
    slot: number;
    item_id: number;
};

export type PlayerNpcRow = {
    npc_id: number;
    character_id: number;
    name: string;
    gender: number;
    skin: number;
    face: number;
    hair: number;
    equips: PlayerNpcEquipRow[];
    map_id: number;
    x: number;
    y: number;
};

const SELECT_COLS = "npc_id, character_id, name, gender, skin, face, hair, equips, map_id, x, y";

export class PlayerNpcRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    async list(worldId: number): Promise<PlayerNpcRow[]> {
        const res = await this.ctx.getPgGlobalPool(worldId).query(
            `SELECT ${SELECT_COLS} FROM player_npcs WHERE world_id = $1 ORDER BY npc_id`,
            [worldId]
        );
        return res.rows as PlayerNpcRow[];
    }

    async insert(worldId: number, row: PlayerNpcRow): Promise<boolean> {
        const res = await this.ctx.getPgGlobalPool(worldId).query(
            `INSERT INTO player_npcs (world_id, npc_id, character_id, name, gender, skin, face, hair, equips, map_id, x, y)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb, $10, $11, $12)
             ON CONFLICT (world_id, npc_id) DO NOTHING`,
            [
                worldId,
                row.npc_id,
                row.character_id,
                row.name,
                row.gender,
                row.skin,
                row.face,
                row.hair,
                JSON.stringify(row.equips),
                row.map_id,
                row.x,
                row.y,
            ]
        );
        return (res.rowCount ?? 0) > 0;
    }

    async remove(worldId: number, npcId: number): Promise<boolean> {
        const res = await this.ctx.getPgGlobalPool(worldId).query(
            "DELETE FROM player_npcs WHERE world_id = $1 AND npc_id = $2",
            [worldId, npcId]
        );
        return (res.rowCount ?? 0) > 0;
    }
}
