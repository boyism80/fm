import type { CharacterSession } from "./session-types";

export function serializeCharacterSession(s: CharacterSession): string {
    return JSON.stringify(s);
}

function isCharacterSessionRecord(o: unknown): o is CharacterSession {
    if (o == null || typeof o !== "object") {
        return false;
    }
    const x = o as Record<string, unknown>;
    if (x.version !== 1) {
        return false;
    }
    if (typeof x.worldId !== "number" || !Number.isFinite(x.worldId)) {
        return false;
    }
    if (typeof x.accountId !== "number" || !Number.isFinite(x.accountId)) {
        return false;
    }
    if (typeof x.characterId !== "number" || !Number.isFinite(x.characterId)) {
        return false;
    }
    if (x.characterName != null && typeof x.characterName !== "string") {
        return false;
    }
    if (typeof x.state !== "string") {
        return false;
    }
    const gs = x.gameServer;
    if (gs == null || typeof gs !== "object") {
        return false;
    }
    const g = gs as Record<string, unknown>;
    if (g.id != null && typeof g.id !== "string") {
        return false;
    }
    if (g.worldId != null && (typeof g.worldId !== "number" || !Number.isFinite(g.worldId))) {
        return false;
    }
    if (g.channelId != null && (typeof g.channelId !== "number" || !Number.isFinite(g.channelId))) {
        return false;
    }
    if (typeof g.connected !== "boolean") {
        return false;
    }
    const ts = x.timestamps;
    if (ts == null || typeof ts !== "object") {
        return false;
    }
    const t = ts as Record<string, unknown>;
    if (t.createdAt != null && typeof t.createdAt !== "string") {
        return false;
    }
    if (t.updatedAt != null && typeof t.updatedAt !== "string") {
        return false;
    }
    return true;
}

export function deserializeCharacterSession(raw: string): CharacterSession | null {
    try {
        const o = JSON.parse(raw) as unknown;
        if (!isCharacterSessionRecord(o)) {
            return null;
        }
        const gs = o.gameServer;
        const ts = o.timestamps;
        return {
            version: o.version,
            worldId: o.worldId,
            accountId: o.accountId,
            characterId: o.characterId,
            characterName: o.characterName ?? null,
            state: o.state,
            gameServer: {
                id: gs.id ?? null,
                worldId: gs.worldId ?? null,
                channelId: gs.channelId ?? null,
                connected: gs.connected,
            },
            timestamps: {
                createdAt: ts.createdAt ?? null,
                updatedAt: ts.updatedAt ?? null,
            },
        };
    } catch {
        return null;
    }
}
