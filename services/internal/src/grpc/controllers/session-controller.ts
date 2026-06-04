import { CHARACTER_MODEL, CHARACTER_PERSISTED } from "../character-persisted";
import { BUDDY_ENTRY, BUDDY_LIST_ENTRY } from "../buddy-persisted";
import { makeKeyLayoutProtoList } from "../key-layout-io";
import { INVENTORY_MODEL, INVENTORY_PERSISTED } from "../inventory-persisted";
import { SKILL_MODEL, SKILL_PERSISTED } from "../skill-persisted";
import { BUFF_MODEL, BUFF_PERSISTED } from "../buff-persisted";
import { grpcMapper } from "../mappers";
import { SessionErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type { BuddyService, BuddyListEntry } from "../../services/buddy-service";
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
import type {
    BuddyEntry,
    CharacterPersisted,
    InventoryPersisted,
    SkillPersisted,
    BuffPersisted,
} from "../../protobuf/generated/fminternal/internal_service";
import type { CharacterModel } from "../../repos/character-repository";
import type { InventoryModel } from "../../repos/inventory-repository";
import type { SkillModel } from "../../repos/skill-repository";
import type { BuffModel } from "../../repos/buff-repository";
import { Controller, Method } from "../grpc-method-decorator";

type OwnedCharacterResult =
    | { kind: "skip" }
    | { kind: "bad"; code: SessionErrorCode }
    | { kind: "ok"; row: CharacterModel };

@Controller("sessionController")
export class SessionGrpcController {
    private readonly characterService: CharacterService;
    private readonly inventoryRepository: InventoryRepository;
    private readonly skillService: SkillService;
    private readonly buffService: BuffService;
    private readonly sessionService: SessionService;
    private readonly buddyService: BuddyService;
    private readonly internalConfig: Pick<InternalConfig, "game_servers">;
    private readonly characterRealtimeStateRepository: CharacterRealtimeStateRepository;
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        characterService: CharacterService,
        inventoryRepository: InventoryRepository,
        skillService: SkillService,
        buffService: BuffService,
        sessionService: SessionService,
        buddyService: BuddyService,
        internalConfig: Pick<InternalConfig, "game_servers">,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        grpcError: GrpcErrorHandler
    ) {
        this.characterService = characterService;
        this.inventoryRepository = inventoryRepository;
        this.skillService = skillService;
        this.buffService = buffService;
        this.sessionService = sessionService;
        this.buddyService = buddyService;
        this.internalConfig = internalConfig;
        this.characterRealtimeStateRepository = characterRealtimeStateRepository;
        this.grpcError = grpcError;
    }

    private async ownedCharacterForSessionRpc(
        worldId: number,
        accountId: number,
        characterId: number | undefined,
    ): Promise<OwnedCharacterResult> {
        if (characterId == null || characterId === 0) {
            return { kind: "skip" };
        }
        const row = await this.characterService.getCharacter(worldId, characterId);
        if (!row || row.accountId !== accountId) {
            return { kind: "bad", code: SessionErrorCode.SESSION_NOT_FOUND };
        }
        return { kind: "ok", row };
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
            const transition = await this.sessionService.beginTransition(worldId, row.accountId, row.characterId, row.name);
            callback(null, {
                ok: transition.ok,
                errorCode: transition.code ?? SessionErrorCode.SESSION_UNKNOWN,
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
                throw Object.assign(new Error(`Unknown world_id for game attach: ${worldId}`), { code: "UNKNOWN_WORLD" });
            }
            const hasChannel = (worldCfg.channels ?? []).some((ch) => ch.channel_id === channelId);
            if (!hasChannel) {
                throw Object.assign(new Error(`Unknown channel_id for world ${worldId}: ${channelId}`), { code: "UNKNOWN_CHANNEL" });
            }

            const row = await this.characterService.getCharacter(worldId, characterId);
            if (!row) {
                callback(null, {
                    found: false,
                    character: undefined,
                    inventory: [],
                    skills: [],
                    buffs: [],
                    keyLayout: [],
                    partyId: undefined,
                    guildId: undefined,
                    buddies: [],
                    buddyCapacity: 0,
                });
                return;
            }

            const attach = await this.sessionService.attachGameSession(
                worldId,
                row.accountId,
                row.characterId,
                row.name,
                channelId
            );
            if (!attach.ok) {
                throw new Error(`attach game session failed: ${attach.code}`);
            }

            const [inventoryList, skillList, buffList, keyLayoutBindings, buddyPack] = await Promise.all([
                this.inventoryRepository.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
                this.skillService.getSkills(worldId, characterId),
                this.buffService.getBuffs(worldId, characterId),
                this.characterService.getKeyLayoutBindings(worldId, characterId),
                this.buddyService.getAll(worldId, characterId),
            ]);
            const realtime = await this.characterRealtimeStateRepository.get(worldId, characterId);
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
                partyId: realtime?.partyId != null ? realtime.partyId : undefined,
                guildId: realtime?.guildId != null ? realtime.guildId : undefined,
                buddies: buddyPack.buddies.map((buddy) =>
                    grpcMapper.map<BuddyListEntry, BuddyEntry>(buddy, BUDDY_LIST_ENTRY, BUDDY_ENTRY)
                ),
                buddyCapacity: buddyPack.capacity >>> 0,
            });
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
            if (oc.kind === "bad") {
                callback(null, {
                    ok: false,
                    errorCode: oc.code,
                });
                return;
            }
            const result = await this.sessionService.refresh(
                call.request.worldId,
                call.request.accountId,
                call.request.characterId
            );
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
            if (oc.kind === "bad") {
                callback(null, { ok: false });
                return;
            }
            const characterName = oc.kind === "ok" ? oc.row.name : undefined;
            const result = await this.sessionService.logout(call.request.worldId, call.request.accountId, {
                disconnectSource: call.request.disconnectSource,
                transferDisconnect: call.request.transferDisconnect,
                characterId: call.request.characterId,
                characterName,
                channelId: call.request.channelId,
            });
            callback(null, { ok: result.ok });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
