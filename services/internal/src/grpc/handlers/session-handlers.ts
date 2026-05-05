import { makeCharacterMessage } from "../character-persisted";
import { makeKeyLayoutProtoList } from "../key-layout-io";
import { makeInventoryMessage } from "../inventory-persisted";
import { makeSkillMessage } from "../skill-persisted";
import { SessionErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type { CharacterService } from "../../services/character-service";
import type { SkillService } from "../../services/skill-service";
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

export function createSessionHandlers(
    characterService: CharacterService,
    inventoryRepository: InventoryRepository,
    skillService: SkillService,
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
                    ok: Boolean(transition.ok),
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

                const [inventoryList, skillList, keyLayoutBindings] = await Promise.all([
                    inventoryRepository.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
                    skillService.getSkills(worldId, characterId),
                    characterService.getKeyLayoutBindings(worldId, characterId),
                ]);
                const realtime = await characterRealtimeStateRepository.get(worldId, characterId);
                const guildId = realtime?.guildId ?? 0;
                callback(null, {
                    found: true,
                    character: makeCharacterMessage(row),
                    inventory: inventoryList.map(makeInventoryMessage),
                    skills: skillList.map(makeSkillMessage),
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
                    ok: Boolean(result.ok),
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
                callback(null, { ok: Boolean(result.ok) });
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}
