export interface AccountSession {
    version: number;
    worldId: number | null;
    accountId: number;
    state: string;
    loginServer: { id: string | null; connected: boolean };
    timestamps: { createdAt: string | null; updatedAt: string | null; stateChangedAt: string | null };
}

export interface CharacterSession {
    version: number;
    worldId: number;
    accountId: number;
    characterId: number;
    characterName: string | null;
    state: string;
    gameServer: { id: string | null; worldId: number | null; channelId: number | null; connected: boolean };
    timestamps: { createdAt: string | null; updatedAt: string | null };
}
