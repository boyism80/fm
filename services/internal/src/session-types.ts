import type { AccountSessionState, CharacterSessionState } from "./session-state";

export interface AccountSession {
    version: number;
    worldId: number | null;
    accountId: number;
    state: AccountSessionState;
    loginServer: { id: string | null; connected: boolean };
    timestamps: { createdAt: string | null; updatedAt: string | null; stateChangedAt: string | null };
}

export interface CharacterSession {
    version: number;
    worldId: number;
    accountId: number;
    characterId: number;
    characterName: string | null;
    state: CharacterSessionState;
    gameServer: { id: string | null; worldId: number | null; channelId: number | null; connected: boolean };
    timestamps: { createdAt: string | null; updatedAt: string | null };
}
