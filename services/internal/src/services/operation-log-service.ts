import type { InternalContext } from "../context/internal-context";

export class OperationLogService {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    async write(worldId: number, channelId: number, characterId: number, kind: string, meso: number, detail: string): Promise<boolean> {
        if (!Number.isFinite(worldId) || worldId < 0) {
            return false;
        }
        if (kind.length <= 0) {
            return false;
        }
        if (kind === "wedding_guest_invite_failure") {
            await this.ctx.recordApplicationFailure(worldId, kind, `character=${characterId} channel=${channelId}`, detail);
            return true;
        }
        await this.ctx.withPgGlobalTransaction(worldId, async (txClient) => {
            await txClient.query(
                "INSERT INTO operation_logs (channel_id, character_id, kind, meso, detail) VALUES ($1, $2, $3, $4, $5)",
                [channelId, characterId, kind, meso, detail],
            );
        });
        return true;
    }
}
