import { SecondPasswordStatus, SessionErrorCode } from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { AccountService } from "../../services/account-service";
import type { CharacterService } from "../../services/character-service";
import type { SessionService } from "../../services/session-service";
import type { InternalConfig } from "@fm/shared/types/internal-config";
import type {
    BeginGameTransitionReply,
    BeginGameTransitionRequest,
    AuthorizeCharacterSelectionReply,
    AuthorizeCharacterSelectionRequest,
    EnterCashShopReply,
    EnterCashShopRequest,
    EnterGameReply,
    EnterGameRequest,
    LogoutSessionReply,
    LogoutSessionRequest,
    RefreshSessionReply,
    RefreshSessionRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import type { CharacterModel } from "@fm/shared/repos/character-repository";
import { Controller, Method } from "../grpc-method-decorator";

enum OwnedCharacterKind {
    Skip = 0,
    Bad = 1,
    Ok = 2,
}

type OwnedCharacterResult =
    | { kind: OwnedCharacterKind.Skip }
    | { kind: OwnedCharacterKind.Bad; code: SessionErrorCode }
    | { kind: OwnedCharacterKind.Ok; row: CharacterModel };

@Controller("sessionController")
export class SessionGrpcController {
    private readonly characterService: CharacterService;
    private readonly sessionService: SessionService;
    private readonly accountService: AccountService;
    private readonly internalConfig: Pick<InternalConfig, "game_servers">;
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        characterService: CharacterService,
        sessionService: SessionService,
        accountService: AccountService,
        internalConfig: Pick<InternalConfig, "game_servers">,
        grpcError: GrpcErrorHandler
    ) {
        this.characterService = characterService;
        this.sessionService = sessionService;
        this.accountService = accountService;
        this.internalConfig = internalConfig;
        this.grpcError = grpcError;
    }

    private async ownedCharacterForSessionRpc(
        worldId: number,
        accountId: number,
        characterId: number | undefined,
    ): Promise<OwnedCharacterResult> {
        if (characterId == null || characterId === 0) {
            return { kind: OwnedCharacterKind.Skip };
        }
        const row = await this.characterService.getCharacter(worldId, characterId);
        if (!row || (accountId !== 0 && row.accountId !== accountId)) {
            return { kind: OwnedCharacterKind.Bad, code: SessionErrorCode.SESSION_NOT_FOUND };
        }
        return { kind: OwnedCharacterKind.Ok, row };
    }

    @Method("beginGameTransition")
    async beginGameTransition(call: GrpcCall<BeginGameTransitionRequest>, callback: GrpcCallback<BeginGameTransitionReply>) {
        try {
            const worldId = call.request.worldId;
            const accountId = call.request.accountId;
            const characterId = call.request.characterId;
            const row = await this.characterService.getCharacter(worldId, characterId);
            if (!row) {
                throw Object.assign(new Error(`character not found: ${characterId}`), { code: "INVALID_CHARACTER_ID" });
            }
            if (row.accountId !== accountId) {
                throw Object.assign(new Error(`account mismatch for character ${characterId}`), { code: "INVALID_PAYLOAD" });
            }
            const transition = await this.sessionService.beginTransition(
                worldId,
                row.accountId,
                row.characterId,
                row.name,
                call.request.clientIp,
                call.request.debuffs ?? [],
                call.request.sourceChannelId ?? null,
                call.request.sourceCashShopId ?? null,
            );
            callback(null, {
                ok: transition.ok,
                errorCode: transition.code ?? SessionErrorCode.SESSION_UNKNOWN,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("authorizeCharacterSelection")
    async authorizeCharacterSelection(call: GrpcCall<AuthorizeCharacterSelectionRequest>, callback: GrpcCallback<AuthorizeCharacterSelectionReply>) {
        try {
            const worldId = call.request.worldId;
            const accountId = call.request.accountId;
            const characterId = call.request.characterId;
            const row = await this.characterService.getCharacter(worldId, characterId);
            if (!row || row.accountId !== accountId) {
                callback(null, { status: SecondPasswordStatus.SECOND_PASSWORD_CHARACTER_NOT_OWNED });
                return;
            }

            const result = await this.accountService.authorizeCharacterSelection(worldId, accountId, call.request.secondPassword);
            let status = SecondPasswordStatus.SECOND_PASSWORD_INVALID;
            switch (result) {
                case "success":
                    status = SecondPasswordStatus.SECOND_PASSWORD_SUCCESS;
                    break;
                case "already_set":
                    status = SecondPasswordStatus.SECOND_PASSWORD_ALREADY_SET;
                    break;
                case "not_set":
                    status = SecondPasswordStatus.SECOND_PASSWORD_NOT_SET;
                    break;
                case "locked":
                    status = SecondPasswordStatus.SECOND_PASSWORD_LOCKED;
                    break;
            }
            if (status !== SecondPasswordStatus.SECOND_PASSWORD_SUCCESS) {
                callback(null, { status });
                return;
            }

            const transition = await this.sessionService.beginTransition(
                worldId,
                accountId,
                characterId,
                row.name,
                call.request.clientIp,
                [],
                null,
                null,
            );
            callback(null, {
                status: transition.ok
                    ? SecondPasswordStatus.SECOND_PASSWORD_SUCCESS
                    : SecondPasswordStatus.SECOND_PASSWORD_SESSION_ERROR,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("enterGame")
    async enterGame(call: GrpcCall<EnterGameRequest>, callback: GrpcCallback<EnterGameReply>) {
        try {
            const worldId = call.request.worldId;
            const characterId = call.request.characterId;
            const channelId = call.request.channelId;
            const worldCfg = this.internalConfig.game_servers?.worlds?.[String(worldId)];
            if (!worldCfg) {
                throw Object.assign(new Error(`Unknown world_id for enter game: ${worldId}`), { code: "UNKNOWN_WORLD" });
            }
            const hasChannel = (worldCfg.channels ?? []).some((ch) => ch.channel_id === channelId);
            if (!hasChannel) {
                throw Object.assign(new Error(`Unknown channel_id for world ${worldId}: ${channelId}`), { code: "UNKNOWN_CHANNEL" });
            }

            const row = await this.characterService.getCharacter(worldId, characterId);
            if (!row) {
                callback(null, {
                    family: undefined,
                    notes: [],
                    found: false,
                    character: undefined,
                    inventory: [],
                    skills: [],
                    buffs: [],
                    quests: [],
                    savedLocations: [],
                    keyLayout: [],
                    skillMacros: [],
                    partyId: undefined,
                    guildId: undefined,
                    buddies: [],
                    buddyCapacity: 0,
                    debuffs: [],
                    cashWishlist: [],
                    rings: [],
                    monsterBook: [],
                    records: [],
                    mount: undefined,
                    account: undefined,
                    unsettledRequests: [],
                });
                return;
            }
            callback(null, await this.sessionService.enterGame(worldId, row, channelId, call.request.clientIp, () => call.cancelled));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("enterCashShop")
    async enterCashShop(call: GrpcCall<EnterCashShopRequest>, callback: GrpcCallback<EnterCashShopReply>) {
        try {
            const { worldId, characterId, cashShopId } = call.request;
            const worldCfg = this.internalConfig.game_servers?.worlds?.[String(worldId)];
            const hasCashShop = (worldCfg?.cash_shops ?? []).some((cs) => cs.cash_shop_id === cashShopId);
            if (hasCashShop === false) {
                throw Object.assign(new Error(`Unknown cash_shop_id for world ${worldId}: ${cashShopId}`), { code: "UNKNOWN_CHANNEL" });
            }

            const row = await this.characterService.getCharacter(worldId, characterId);
            if (!row) {
                throw Object.assign(new Error(`character not found: ${characterId}`), { code: "INVALID_CHARACTER_ID" });
            }
            callback(null, await this.sessionService.enterCashShop(worldId, row, cashShopId, call.request.clientIp, () => call.cancelled));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("refreshSession")
    async refreshSession(call: GrpcCall<RefreshSessionRequest>, callback: GrpcCallback<RefreshSessionReply>) {
        try {
            const oc = await this.ownedCharacterForSessionRpc(
                call.request.worldId,
                call.request.accountId,
                call.request.characterId,
            );
            if (oc.kind === OwnedCharacterKind.Bad) {
                callback(null, {
                    ok: false,
                    errorCode: oc.code,
                });
                return;
            }
            const { characterId, channelId, cashShopId } = call.request;
            let owner: { characterId: number; channelId?: number; cashShopId?: number } | null = null;
            if (characterId != null && cashShopId != null) {
                owner = { characterId, cashShopId };
            } else if (characterId != null && channelId != null) {
                owner = { characterId, channelId };
            }
            const result = await this.sessionService.refresh(call.request.worldId, call.request.accountId, owner);
            callback(null, {
                ok: result.ok,
                errorCode: result.code ?? SessionErrorCode.SESSION_UNKNOWN,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("logoutSession")
    async logoutSession(call: GrpcCall<LogoutSessionRequest>, callback: GrpcCallback<LogoutSessionReply>) {
        try {
            const oc = await this.ownedCharacterForSessionRpc(
                call.request.worldId,
                call.request.accountId,
                call.request.characterId,
            );
            if (oc.kind === OwnedCharacterKind.Bad) {
                callback(null, { ok: false });
                return;
            }
            const accountId = oc.kind === OwnedCharacterKind.Ok ? oc.row.accountId : call.request.accountId;
            if (accountId === 0) {
                callback(null, { ok: false });
                return;
            }
            const result = await this.sessionService.logout(call.request.worldId, accountId, {
                disconnectSource: call.request.disconnectSource,
                transferDisconnect: call.request.transferDisconnect,
                characterId: call.request.characterId,
                channelId: call.request.channelId,
                cashShopId: call.request.cashShopId,
            });
            callback(null, { ok: result.ok });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
