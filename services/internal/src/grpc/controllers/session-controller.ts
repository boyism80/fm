import { CHARACTER_MODEL, CHARACTER_PROTO } from "../character-proto";
import { BUDDY_ENTRY, BUDDY_LIST_ENTRY } from "../buddy-proto";
import { makeKeyLayoutProtoList } from "../key-layout-io";
import { INVENTORY_MODEL, INVENTORY_PROTO } from "../inventory-proto";
import { SKILL_MODEL, SKILL_PROTO } from "../skill-proto";
import { BUFF_MODEL, BUFF_PROTO } from "../buff-proto";
import { QUEST_MODEL, QUEST_PROTO } from "../quest-proto";
import { SAVED_LOCATION_MODEL, SAVED_LOCATION_PROTO } from "../saved-location-proto";
import { grpcMapper } from "../mappers";
import { SessionErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type { BuddyService, BuddyListEntry } from "../../services/buddy-service";
import type { CharacterService } from "../../services/character-service";
import type { BuffService } from "../../services/buff-service";
import type { SessionService } from "../../services/session-service";
import type { CashShopService } from "../../services/cash-shop-service";
import type { MarriageService } from "../../services/marriage-service";
import type { InventoryRepository } from "../../repos/inventory-repository";
import type { SkillRepository } from "../../repos/skill-repository";
import type { QuestRepository } from "../../repos/quest-repository";
import type { SavedLocationRepository } from "../../repos/saved-location-repository";
import type { CharacterRealtimeStateRepository } from "../../repos/character-realtime-state-repository";
import type { InternalConfig } from "../../types/internal-config";
import type {
    BeginGameTransitionReply,
    BeginGameTransitionRequest,
    EnterCashShopReply,
    EnterCashShopRequest,
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
    Character,
    Inventory,
    Skill,
    Buff,
    Quest,
    SavedLocation,
} from "../../protobuf/generated/fminternal/internal_service";
import type { CharacterModel } from "../../repos/character-repository";
import type { InventoryModel } from "../../repos/inventory-repository";
import type { SkillModel } from "../../repos/skill-repository";
import type { BuffModel } from "../../repos/buff-repository";
import type { QuestModel } from "../../repos/quest-repository";
import type { SavedLocationModel } from "../../repos/saved-location-repository";
import { Controller, Method } from "../grpc-method-decorator";

type OwnedCharacterResult =
    | { kind: "skip" }
    | { kind: "bad"; code: SessionErrorCode }
    | { kind: "ok"; row: CharacterModel };

@Controller("sessionController")
export class SessionGrpcController {
    private readonly characterService: CharacterService;
    private readonly inventoryRepository: InventoryRepository;
    private readonly skillRepository: SkillRepository;
    private readonly buffService: BuffService;
    private readonly questRepository: QuestRepository;
    private readonly savedLocationRepository: SavedLocationRepository;
    private readonly sessionService: SessionService;
    private readonly buddyService: BuddyService;
    private readonly internalConfig: Pick<InternalConfig, "game_servers">;
    private readonly characterRealtimeStateRepository: CharacterRealtimeStateRepository;
    private readonly cashShopService: CashShopService;
    private readonly marriageService: MarriageService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        characterService: CharacterService,
        inventoryRepository: InventoryRepository,
        skillRepository: SkillRepository,
        buffService: BuffService,
        questRepository: QuestRepository,
        savedLocationRepository: SavedLocationRepository,
        sessionService: SessionService,
        buddyService: BuddyService,
        internalConfig: Pick<InternalConfig, "game_servers">,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        cashShopService: CashShopService,
        marriageService: MarriageService,
        grpcError: GrpcErrorHandler
    ) {
        this.characterService = characterService;
        this.inventoryRepository = inventoryRepository;
        this.skillRepository = skillRepository;
        this.buffService = buffService;
        this.questRepository = questRepository;
        this.savedLocationRepository = savedLocationRepository;
        this.sessionService = sessionService;
        this.buddyService = buddyService;
        this.internalConfig = internalConfig;
        this.characterRealtimeStateRepository = characterRealtimeStateRepository;
        this.cashShopService = cashShopService;
        this.marriageService = marriageService;
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
        if (!row || (accountId !== 0 && row.accountId !== accountId)) {
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

    private async loadCharacterReply(worldId: number, row: CharacterModel): Promise<EnterGameReply> {
        const characterId = row.characterId;
        const [inventoryList, skillList, buffList, questList, savedLocationList, keyLayoutBindings, buddyPack, marriage, cashWishlist] = await Promise.all([
            this.inventoryRepository.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.skillRepository.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.buffService.getBuffs(worldId, characterId),
            this.questRepository.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.savedLocationRepository.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.characterService.getKeyLayoutBindings(worldId, characterId),
            this.buddyService.getAll(worldId, characterId),
            this.marriageService.loadMarriage(worldId, characterId),
            this.cashShopService.getWishlist(worldId, characterId),
        ]);
        const realtime = await this.characterRealtimeStateRepository.get(worldId, characterId);
        return {
            found: true,
            character: grpcMapper.map<CharacterModel, Character>(
                row,
                CHARACTER_MODEL,
                CHARACTER_PROTO
            ),
            inventory: inventoryList.map((inventory) =>
                grpcMapper.map<InventoryModel, Inventory>(
                    inventory,
                    INVENTORY_MODEL,
                    INVENTORY_PROTO
                )
            ),
            skills: skillList.map((skill) =>
                grpcMapper.map<SkillModel, Skill>(
                    skill,
                    SKILL_MODEL,
                    SKILL_PROTO
                )
            ),
            buffs: buffList.map((buff) =>
                grpcMapper.map<BuffModel, Buff>(
                    buff,
                    BUFF_MODEL,
                    BUFF_PROTO
                )
            ),
            quests: questList.map((quest) =>
                grpcMapper.map<QuestModel, Quest>(
                    quest,
                    QUEST_MODEL,
                    QUEST_PROTO
                )
            ),
            savedLocations: savedLocationList.map((loc) =>
                grpcMapper.map<SavedLocationModel, SavedLocation>(
                    loc,
                    SAVED_LOCATION_MODEL,
                    SAVED_LOCATION_PROTO
                )
            ),
            keyLayout: makeKeyLayoutProtoList(keyLayoutBindings),
            partyId: realtime?.partyId != null ? realtime.partyId : undefined,
            guildId: realtime?.guildId != null ? realtime.guildId : undefined,
            buddies: buddyPack.buddies.map((buddy) =>
                grpcMapper.map<BuddyListEntry, BuddyEntry>(buddy, BUDDY_LIST_ENTRY, BUDDY_ENTRY)
            ),
            buddyCapacity: buddyPack.capacity >>> 0,
            debuffs: [],
            marriage: marriage.marriage,
            cashWishlist: cashWishlist.filter((sn) => sn !== 0),
        };
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
                    found: false,
                    character: undefined,
                    inventory: [],
                    skills: [],
                    buffs: [],
                    quests: [],
                    savedLocations: [],
                    keyLayout: [],
                    partyId: undefined,
                    guildId: undefined,
                    buddies: [],
                    buddyCapacity: 0,
                    debuffs: [],
                    cashWishlist: [],
                });
                return;
            }
            const reply = await this.loadCharacterReply(worldId, row);

            // Enter the game last: the game server treats an error reply as "did not enter the game".
            const enter = await this.sessionService.enterGame(worldId, row.accountId, row.characterId, channelId, call.request.clientIp);
            if (!enter.ok) {
                throw new Error(`enter game session failed: ${enter.code}`);
            }
            reply.debuffs = enter.debuffs ?? [];
            callback(null, reply);
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
            const game = await this.loadCharacterReply(worldId, row);
            const account = await this.cashShopService.loadAccount(worldId, row.accountId, characterId);

            const enter = await this.sessionService.enterCashShop(worldId, row.accountId, characterId, cashShopId, call.request.clientIp);
            if (!enter.ok) {
                throw new Error(`enter cash shop session failed: ${enter.code}`);
            }
            game.debuffs = enter.debuffs ?? [];
            callback(null, { game, returnChannelId: enter.returnChannelId ?? 0, ...account });
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
            if (oc.kind === "bad") {
                callback(null, { ok: false });
                return;
            }
            const accountId = oc.kind === "ok" ? oc.row.accountId : call.request.accountId;
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
