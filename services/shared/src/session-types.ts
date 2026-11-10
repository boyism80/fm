import type { AccountSessionState } from "./session-state";
import type { CharacterPresentation, Debuff } from "./protobuf/generated/fminternal/internal_service";

export interface AccountSession {
    version: number;
    worldId: number | null;
    accountId: number;
    state: AccountSessionState;
    loginServer: { id: string | null; connected: boolean };
    characterId: number | null;
    characterName: string | null;
    channelId: number | null;
    cashShopId: number | null;
    returnChannelId: number | null;
    gameToGameTransfer: boolean;
    debuffs: Debuff[];
    presentation: CharacterPresentation | null;
    timestamps: {
        createdAt: string | null;
        updatedAt: string | null;
        stateChangedAt: string | null;
        // When the character entered the game from the login server; a move between servers keeps it.
        loggedInAt: string | null;
    };
}
