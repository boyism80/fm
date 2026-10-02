import type { AccountSessionState, CharacterSessionState } from "./session-state";

export interface AccountSession {
    version: number;
    worldId: number | null;
    accountId: number;
    state: AccountSessionState;
    loginServer: { id: string | null; connected: boolean };
    characterId: number | null;
    characterName: string | null;
    gameServer: { channelId: number | null; connected: boolean };
    gameToGameTransfer: boolean;
    timestamps: { createdAt: string | null; updatedAt: string | null; stateChangedAt: string | null };
}

export interface CharacterSession {
    worldId: number;
    accountId: number;
    characterId: number;
    characterName: string | null;
    state: CharacterSessionState;
    gameServer: { channelId: number | null; connected: boolean };
}
