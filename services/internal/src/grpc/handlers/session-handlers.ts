import { CHARACTER_MODEL, CHARACTER_PERSISTED } from "../character-persisted";
import { makeKeyLayoutProtoList } from "../key-layout-io";
import { INVENTORY_MODEL, INVENTORY_PERSISTED } from "../inventory-persisted";
import { SKILL_MODEL, SKILL_PERSISTED } from "../skill-persisted";
import { BUFF_MODEL, BUFF_PERSISTED } from "../buff-persisted";
import { grpcMapper } from "../mappers";
import { SessionErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type { CharacterService } from "../../services/character-service";
import type { SkillService } from "../../services/skill-service";
import type { BuffService } from "../../services/buff-service";
import type { SessionService } from "../../services/session-service";
import type { InventoryRepository } from "../../repos/inventory-repository";
import type { CharacterRealtimeStateRepository } from "../../repos/character-realtime-state-repository";
import type { InternalConfig } from "../../types/internal-config";
import type {
    BeginGameTransitionReply,
    BeginGameTransitionRequest,
    EnterGameReply,
    EnterGameRequest,
    LogoutSessionReply,
    LogoutSessionRequest,
    RefreshSessionReply,
    RefreshSessionRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import type { CharacterPersisted, InventoryPersisted, SkillPersisted, BuffPersisted } from "../../protobuf/generated/fminternal/internal_service";
import type { CharacterModel } from "../../repos/character-repository";
import type { InventoryModel } from "../../repos/inventory-repository";
import type { SkillModel } from "../../repos/skill-repository";
import type { BuffModel } from "../../repos/buff-repository";
import { GrpcController, GrpcMethod } from "../grpc-method-decorator";

export function createSessionHandlers(
    characterService: CharacterService,
    inventoryRepository: InventoryRepository,
    skillService: SkillService,
    buffService: BuffService,
    sessionService: SessionService,
    internalConfig: Pick<InternalConfig, "game_servers">,
    characterRealtimeStateRepository: CharacterRealtimeStateRepository,
    grpcError: GrpcErrorHandler
) {
    return {
        async beginGameTransition(call: GrpcCall<BeginGameTransitionRequest>, callback: GrpcCallback<BeginGameTransitionReply>) {
            try {
                const worldId = call.request.worldId;
                const accountId = call.request.accountId;
                const characterId = call.request.characterId;
                const row = await characterService.getCharacter(worldId, characterId);
                if (!row) {
                    throw Object.assign(new Error(`character not found: ${characterId}`), { code: "INVALID_CHARACTER_ID" });
                }
                if (row.accountId !== accountId) {
                    throw Object.assign(new Error(`account mismatch for character ${characterId}`), { code: "INVALID_PAYLOAD" });
                }
                const transition = await sessionService.beginTransition(worldId, row.accountId, row.characterId, row.name);
                callback(null, {
                    ok: transition.ok,
                    errorCode: transition.code ?? SessionErrorCode.SESSION_UNKNOWN,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async enterGame(call: GrpcCall<EnterGameRequest>, callback: GrpcCallback<EnterGameReply>) {
            try {
                const worldId = call.request.worldId;
                const characterId = call.request.characterId;
                const channelId = call.request.channelId;
                const worldCfg = internalConfig.game_servers?.worlds?.[String(worldId)];
                if (!worldCfg) {
                    throw Object.assign(new Error(`Unknown world_id for game attach: ${worldId}`), { code: "UNKNOWN_WORLD" });
                }
                const hasChannel = (worldCfg.channels ?? []).some((ch) => ch.channel_id === channelId);
                if (!hasChannel) {
                    throw Object.assign(new Error(`Unknown channel_id for world ${worldId}: ${channelId}`), { code: "UNKNOWN_CHANNEL" });
                }

                const row = await characterService.getCharacter(worldId, characterId);
                if (!row) {
                    callback(null, {
                        found: false,
                        character: undefined,
                        inventory: [],
                        skills: [],
                        buffs: [],
                        keyLayout: [],
                        partyId: undefined,
                        guildId: 0,
                    });
                    return;
                }

                const attach = await sessionService.attachGameSession(worldId, row.accountId, channelId);
                if (!attach.ok) {
                    throw new Error(`attach game session failed: ${attach.code}`);
                }

                const [inventoryList, skillList, buffList, keyLayoutBindings] = await Promise.all([
                    inventoryRepository.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
                    skillService.getSkills(worldId, characterId),
                    buffService.getBuffs(worldId, characterId),
                    characterService.getKeyLayoutBindings(worldId, characterId),
                ]);
                const realtime = await characterRealtimeStateRepository.get(worldId, characterId);
                const guildId = realtime?.guildId ?? 0;
                callback(null, {
                    found: true,
                    character: grpcMapper.map<CharacterModel, CharacterPersisted>(
                        row,
                        CHARACTER_MODEL,
                        CHARACTER_PERSISTED
                    ),
                    inventory: inventoryList.map((inventory) =>
                        grpcMapper.map<InventoryModel, InventoryPersisted>(
                            inventory,
                            INVENTORY_MODEL,
                            INVENTORY_PERSISTED
                        )
                    ),
                    skills: skillList.map((skill) =>
                        grpcMapper.map<SkillModel, SkillPersisted>(
                            skill,
                            SKILL_MODEL,
                            SKILL_PERSISTED
                        )
                    ),
                    buffs: buffList.map((buff) =>
                        grpcMapper.map<BuffModel, BuffPersisted>(
                            buff,
                            BUFF_MODEL,
                            BUFF_PERSISTED
                        )
                    ),
                    keyLayout: makeKeyLayoutProtoList(keyLayoutBindings),
                    partyId: realtime != null ? realtime.partyId ?? undefined : undefined,
                    guildId,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async refreshSession(call: GrpcCall<RefreshSessionRequest>, callback: GrpcCallback<RefreshSessionReply>) {
            try {
                const result = await sessionService.refresh(call.request.worldId, call.request.accountId);
                callback(null, {
                    ok: result.ok,
                    errorCode: result.code ?? SessionErrorCode.SESSION_UNKNOWN,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async logoutSession(call: GrpcCall<LogoutSessionRequest>, callback: GrpcCallback<LogoutSessionReply>) {
            try {
                const result = await sessionService.logout(call.request.worldId, call.request.accountId, {
                    disconnectSource: call.request.disconnectSource,
                    transferDisconnect: call.request.transferDisconnect,
                });
                callback(null, { ok: result.ok });
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

@GrpcController("sessionController")
export class SessionGrpcController {
    private readonly handlers: ReturnType<typeof createSessionHandlers>;

    constructor(
        characterService: CharacterService,
        inventoryRepository: InventoryRepository,
        skillService: SkillService,
        buffService: BuffService,
        sessionService: SessionService,
        internalConfig: Pick<InternalConfig, "game_servers">,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        grpcError: GrpcErrorHandler
    ) {
        this.handlers = createSessionHandlers(
            characterService,
            inventoryRepository,
            skillService,
            buffService,
            sessionService,
            internalConfig,
            characterRealtimeStateRepository,
            grpcError
        );
    }

    @GrpcMethod("beginGameTransition")
    async beginGameTransition(call: GrpcCall<BeginGameTransitionRequest>, callback: GrpcCallback<BeginGameTransitionReply>) {
        return this.handlers.beginGameTransition(call, callback);
    }

    @GrpcMethod("enterGame")
    async enterGame(call: GrpcCall<EnterGameRequest>, callback: GrpcCallback<EnterGameReply>) {
        return this.handlers.enterGame(call, callback);
    }

    @GrpcMethod("refreshSession")
    async refreshSession(call: GrpcCall<RefreshSessionRequest>, callback: GrpcCallback<RefreshSessionReply>) {
        return this.handlers.refreshSession(call, callback);
    }

    @GrpcMethod("logoutSession")
    async logoutSession(call: GrpcCall<LogoutSessionRequest>, callback: GrpcCallback<LogoutSessionReply>) {
        return this.handlers.logoutSession(call, callback);
    }
}
