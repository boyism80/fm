import { isValidCharacterSessionState } from "./session-state";
import type { CharacterSession } from "./session-types";

export function serializeCharacterSession(s: CharacterSession): string {
    return JSON.stringify(s);
}

type CharacterSessionJson = {
    version: number;
    worldId: number;
    accountId: number;
    characterId: number;
    characterName: string | null;
    state: number;
    gameServer: {
        id: string | null;
        worldId: number | null;
        channelId: number | null;
        connected: boolean;
    };
    timestamps: {
        createdAt: string | null;
        updatedAt: string | null;
    };
};

function isCharacterSessionRecord(o: CharacterSessionJson | null): o is CharacterSessionJson {
    if (o == null) {
        return false;
    }
    if (o.version !== 1) {
        return false;
    }
    if (!Number.isFinite(o.worldId)) {
        return false;
    }
    if (!Number.isFinite(o.accountId)) {
        return false;
    }
    if (!Number.isFinite(o.characterId)) {
        return false;
    }
    if (o.characterName != null && typeof o.characterName !== "string") {
        return false;
    }
    if (!isValidCharacterSessionState(o.state)) {
        return false;
    }
    const gs = o.gameServer;
    if (gs.id != null && typeof gs.id !== "string") {
        return false;
    }
    if (gs.worldId != null && !Number.isFinite(gs.worldId)) {
        return false;
    }
    if (gs.channelId != null && !Number.isFinite(gs.channelId)) {
        return false;
    }
    if (typeof gs.connected !== "boolean") {
        return false;
    }
    const ts = o.timestamps;
    if (ts.createdAt != null && typeof ts.createdAt !== "string") {
        return false;
    }
    if (ts.updatedAt != null && typeof ts.updatedAt !== "string") {
        return false;
    }
    return true;
}

export function deserializeCharacterSession(raw: string): CharacterSession | null {
    try {
        const o = JSON.parse(raw) as CharacterSessionJson;
        if (!isCharacterSessionRecord(o)) {
            return null;
        }
        return {
            version: o.version,
            worldId: o.worldId,
            accountId: o.accountId,
            characterId: o.characterId,
            characterName: o.characterName ?? null,
            state: o.state,
            gameServer: {
                id: o.gameServer.id ?? null,
                worldId: o.gameServer.worldId ?? null,
                channelId: o.gameServer.channelId ?? null,
                connected: o.gameServer.connected,
            },
            timestamps: {
                createdAt: o.timestamps.createdAt ?? null,
                updatedAt: o.timestamps.updatedAt ?? null,
            },
        };
    } catch {
        return null;
    }
}
