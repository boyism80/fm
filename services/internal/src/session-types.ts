import type { AccountSessionState } from "./session-state";
import type { Debuff } from "./protobuf/generated/fminternal/internal_service";

export interface AccountSession {
    version: number;
    worldId: number | null;
    accountId: number;
    state: AccountSessionState;
    loginServer: { id: string | null; connected: boolean };
    characterId: number | null;
    characterName: string | null;
    channelId: number | null;
    gameToGameTransfer: boolean;
    debuffs: Debuff[];
    timestamps: { createdAt: string | null; updatedAt: string | null; stateChangedAt: string | null };
}
