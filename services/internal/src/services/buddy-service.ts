import type { PoolClient } from "pg";
import { BuddyErrorCode, CharacterSessionState } from "../protobuf/generated/fminternal/internal_service";
import type { InternalContext } from "../context/internal-context";
import type { AppConfiguration } from "../config/app-configuration";
import { CharacterBuddyRepository } from "../repos/character-buddy-repository";
import type { CharacterBuddyModel } from "../repos/character-buddy-repository";
import { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import { CharacterRepository } from "../repos/character-repository";
import { SessionRepository } from "../repos/session-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { RabbitMQService } from "./rabbitmq-service";

const DEFAULT_BUDDY_GROUP = "그룹 미지정";
const AMQ_DIRECT_EXCHANGE = "amq.direct";
const BUDDY_EVT_CHANNEL_UPDATE = "channel_update";
const BUDDY_EVT_LIST_UPDATE = "list_update";
const BUDDY_EVT_ADD_REQUEST = "add_request";
const BUDDY_EVT_MULTI_CHAT = "multi_chat";
const BUDDY_SYNC_ACTION_UPDATE = 10;
const MAX_MULTI_CHAT_MSG_LEN = 500;
const MAX_MULTI_CHAT_NAME_LEN = 32;
const MAX_BUDDY_NAME_LEN = 13;
const MAX_BUDDY_GROUP_LEN = 16;

export type BuddyMutationResult = { ok: boolean; code?: BuddyErrorCode };

export type BuddyListEntry = {
    characterId: number;
    name: string;
    groupName: string;
    pending: boolean;
    channelIndex: number;
};

export type RequestBuddyResult = BuddyMutationResult & {
    targetCharacterId?: number;
    targetChannelId?: number;
    requesterView?: BuddyListEntry;
};

export type AcceptBuddyResult = BuddyMutationResult & {
    requesterCharacterId?: number;
    requesterChannelId?: number;
    accepterView?: BuddyListEntry;
    requesterView?: BuddyListEntry;
};

export type RemoveBuddyResult = BuddyMutationResult & {
    buddyCharacterId?: number;
    buddyChannelId?: number;
    buddyWasAccepted?: boolean;
    removedFromOwner?: boolean;
};

export type BroadcastBuddyMultiChatResult = { ok: boolean; code?: BuddyErrorCode; deliveredCount?: number };

export class BuddyService {
    private readonly ctx: InternalContext;
    private readonly app: AppConfiguration;
    private readonly buddyRepo: CharacterBuddyRepository;
    private readonly realtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly characterRepo: CharacterRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;

    constructor(
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        characterBuddyRepository: CharacterBuddyRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        characterRepository: CharacterRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService
    ) {
        this.ctx = internalContext;
        this.app = appConfiguration;
        this.buddyRepo = characterBuddyRepository;
        this.realtimeStateRepo = characterRealtimeStateRepository;
        this.characterRepo = characterRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
    }

    async getAll(worldId: number, characterId: number) {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        const capacity = await this.realtimeStateRepo.getBuddyCapacity(worldId, characterId);
        const all = await this.buddyRepo.getAll(worldId, String(characterId));
        const buddies: BuddyListEntry[] = [];
        for (const model of all.values()) {
            buddies.push(await this.buildEntry(worldId, model));
        }
        buddies.sort((a, b) => a.characterId - b.characterId);
        return { capacity, buddies };
    }

    async applyBuddyChannelIndex(worldId: number, characterId: number, channelIndex: number) {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        const ch = channelIndex;
        if (!Number.isInteger(ch) || ch < -1) {
            return;
        }
        const notifyCharacterIds = await this.buddyRepo.listAcceptedBuddyCharacterIds(worldId, characterId);
        await this.publishBuddyChannelUpdate(worldId, characterId, ch, notifyCharacterIds);
    }

    async requestBuddy(
        worldId: number,
        requesterCharacterId: number,
        targetCharacterName: string,
        groupName: string
    ): Promise<RequestBuddyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(requesterCharacterId);
        const trimmedName = targetCharacterName.trim();
        const trimmedGroup = groupName.trim();
        if (trimmedName.length > MAX_BUDDY_NAME_LEN || trimmedGroup.length > MAX_BUDDY_GROUP_LEN) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_UNKNOWN };
        }

        const requesterRow = await this.characterRepo.get(worldId, requesterCharacterId);
        if (!requesterRow) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_UNKNOWN };
        }

        const targetEntry = await this.unifiedRepo.findCharacterNameEntry(trimmedName);
        if (!targetEntry || targetEntry.world_id !== worldId) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_CHARACTER_NOT_FOUND };
        }
        const targetCharacterId = targetEntry.character_id;
        if (targetCharacterId === requesterCharacterId) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_CANNOT_ADD_SELF };
        }

        const targetRow = await this.characterRepo.get(worldId, targetCharacterId);
        if (!targetRow) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_CHARACTER_NOT_FOUND };
        }

        const existingEntry = await this.buddyRepo.getItem(worldId, String(requesterCharacterId), targetCharacterId);
        if (existingEntry) {
            if (existingEntry.groupName === trimmedGroup) {
                return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_ALREADY_ON_LIST };
            }
            const updated = await this.buddyRepo.set(worldId, {
                ...existingEntry,
                groupName: trimmedGroup,
            });
            const requesterView = await this.buildEntry(worldId, updated);
            const targetChannelIndex = await this.sessionChannelIndex(worldId, targetCharacterId);
            return {
                ok: true,
                targetCharacterId,
                targetChannelId: targetChannelIndex < 0 ? 0 : targetChannelIndex,
                requesterView,
            };
        }

        const requesterCapacity = await this.realtimeStateRepo.getBuddyCapacity(worldId, requesterCharacterId);
        const requesterCount = await this.buddyRepo.countAll(worldId, requesterCharacterId);
        if (requesterCount >= requesterCapacity) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_LIST_FULL };
        }

        const targetAccepted = await this.buddyRepo.countAccepted(worldId, targetCharacterId);
        const targetCapacity = await this.realtimeStateRepo.getBuddyCapacity(worldId, targetCharacterId);
        if (targetAccepted >= targetCapacity) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_TARGET_LIST_FULL };
        }

        const reverse = await this.buddyRepo.getItem(worldId, String(targetCharacterId), requesterCharacterId);
        if (reverse && !reverse.pending) {
            const requesterBuddy = await this.buddyRepo.set(worldId, {
                characterId: requesterCharacterId,
                buddyCharacterId: targetCharacterId,
                groupName: trimmedGroup,
                pending: false,
            });
            await this.buddyRepo.evictGroupCache(worldId, String(requesterCharacterId));
            const requesterChannelIndex = await this.sessionChannelIndex(worldId, requesterCharacterId);
            await this.publishBuddyChannelUpdate(worldId, requesterCharacterId, requesterChannelIndex, [
                targetCharacterId,
            ]);
            const targetChannelIndex = await this.sessionChannelIndex(worldId, targetCharacterId);
            const requesterView = await this.buildEntry(worldId, requesterBuddy);
            return {
                ok: true,
                targetCharacterId,
                targetChannelId: targetChannelIndex < 0 ? 0 : targetChannelIndex,
                requesterView,
            };
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            if (!reverse) {
                await this.buddyRepo.set(
                    worldId,
                    {
                        characterId: targetCharacterId,
                        buddyCharacterId: requesterCharacterId,
                        groupName: DEFAULT_BUDDY_GROUP,
                        pending: true,
                    },
                    { txClient }
                );
            }
            const requesterBuddy = await this.buddyRepo.set(
                worldId,
                {
                    characterId: requesterCharacterId,
                    buddyCharacterId: targetCharacterId,
                    groupName: trimmedGroup,
                    pending: true,
                },
                { txClient }
            );
            return { requesterBuddy };
        });

        await this.buddyRepo.evictGroupCache(worldId, String(requesterCharacterId));
        await this.buddyRepo.evictGroupCache(worldId, String(targetCharacterId));

        const targetChannelIndex = await this.sessionChannelIndex(worldId, targetCharacterId);
        const requesterView = await this.buildEntry(worldId, result.requesterBuddy);
        const targetIncoming = await this.buddyRepo.getItem(worldId, String(targetCharacterId), requesterCharacterId);
        if (targetIncoming) {
            const targetView = await this.buildEntry(worldId, targetIncoming);
            await this.publishBuddyListUpdate(worldId, [targetCharacterId], BUDDY_SYNC_ACTION_UPDATE, [targetView]);
            await this.publishBuddyAddRequest(
                worldId,
                targetCharacterId,
                requesterCharacterId,
                requesterRow.name
            );
        }
        return {
            ok: true,
            targetCharacterId,
            targetChannelId: targetChannelIndex < 0 ? 0 : targetChannelIndex,
            requesterView,
        };
    }

    async acceptBuddy(worldId: number, accepterCharacterId: number, requesterCharacterId: number): Promise<AcceptBuddyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(accepterCharacterId);
        this.assertCharacterId(requesterCharacterId);

        const incoming = await this.buddyRepo.getItem(worldId, String(accepterCharacterId), requesterCharacterId);
        if (!incoming) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_NOT_PENDING };
        }
        if (!incoming.pending) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_ALREADY_ON_LIST };
        }

        const accepterCapacity = await this.realtimeStateRepo.getBuddyCapacity(worldId, accepterCharacterId);
        const accepterAccepted = await this.buddyRepo.countAccepted(worldId, accepterCharacterId);
        if (accepterAccepted >= accepterCapacity) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_LIST_FULL };
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const accepterBuddy = await this.buddyRepo.set(
                worldId,
                {
                    characterId: accepterCharacterId,
                    buddyCharacterId: requesterCharacterId,
                    groupName: DEFAULT_BUDDY_GROUP,
                    pending: false,
                },
                { txClient }
            );
            const onRequester = await this.buddyRepo.getItem(
                worldId,
                String(requesterCharacterId),
                accepterCharacterId,
                { txClient }
            );
            const requesterBuddy = await this.buddyRepo.set(
                worldId,
                onRequester
                    ? { ...onRequester, pending: false }
                    : {
                          characterId: requesterCharacterId,
                          buddyCharacterId: accepterCharacterId,
                          groupName: DEFAULT_BUDDY_GROUP,
                          pending: false,
                      },
                { txClient }
            );
            return { accepterBuddy, requesterBuddy };
        });

        await this.buddyRepo.evictGroupCache(worldId, String(accepterCharacterId));
        await this.buddyRepo.evictGroupCache(worldId, String(requesterCharacterId));

        const requesterChannelIndex = await this.sessionChannelIndex(worldId, requesterCharacterId);
        const accepterChannelIndex = await this.sessionChannelIndex(worldId, accepterCharacterId);
        await this.publishBuddyChannelUpdate(worldId, accepterCharacterId, accepterChannelIndex, [requesterCharacterId]);

        const accepterView = await this.buildEntry(worldId, result.accepterBuddy);
        const requesterView = await this.buildEntry(worldId, result.requesterBuddy);
        await this.publishBuddyListUpdate(worldId, [requesterCharacterId], BUDDY_SYNC_ACTION_UPDATE, [requesterView]);
        return {
            ok: true,
            requesterCharacterId,
            requesterChannelId: requesterChannelIndex < 0 ? 0 : requesterChannelIndex,
            accepterView,
            requesterView,
        };
    }

    async removeBuddy(worldId: number, characterId: number, buddyCharacterId: number): Promise<RemoveBuddyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        this.assertCharacterId(buddyCharacterId);

        const row = await this.buddyRepo.getItem(worldId, String(characterId), buddyCharacterId);
        if (!row) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_UNKNOWN, removedFromOwner: false };
        }

        const buddyWasAccepted = !row.pending;
        await this.buddyRepo.del(worldId, String(characterId), String(buddyCharacterId));
        await this.buddyRepo.evictGroupCache(worldId, String(characterId));

        if (buddyWasAccepted) {
            await this.publishBuddyChannelUpdate(worldId, characterId, -1, [buddyCharacterId]);
        }

        const buddyChannelIndex = await this.sessionChannelIndex(worldId, buddyCharacterId);
        return {
            ok: true,
            buddyCharacterId,
            buddyChannelId: buddyChannelIndex < 0 ? 0 : buddyChannelIndex,
            buddyWasAccepted,
            removedFromOwner: true,
        };
    }

    private assertWorld(worldId: number) {
        const wid = String(worldId);
        if (!this.app.postgresql.worlds[wid]) {
            const err = new Error(`Unknown world_id: ${worldId}`) as Error & { code?: string };
            err.code = "UNKNOWN_WORLD";
            throw err;
        }
    }

    private assertCharacterId(characterId: number) {
        const n = characterId;
        if (!Number.isInteger(n) || n <= 0 || n > 0xffffffff) {
            const err = new Error("character_id must be a positive uint32") as Error & { code?: string };
            err.code = "INVALID_CHARACTER_ID";
            throw err;
        }
    }

    private async sessionChannelIndex(worldId: number, characterId: number) {
        const row = await this.characterRepo.get(worldId, characterId);
        if (!row) {
            return -1;
        }
        const sess = await this.sessionRepo.getCharacterSessionByName(worldId, row.name);
        if (sess?.state !== CharacterSessionState.CHARACTER_SESSION_STATE_ONLINE || sess?.gameServer?.connected !== true) {
            return -1;
        }
        const ch = sess.gameServer?.channelId;
        if (ch == null || !Number.isFinite(ch) || ch < 0) {
            return -1;
        }
        return ch;
    }

    private async publishBuddyChannelUpdate(
        worldId: number,
        characterId: number,
        channelIndex: number,
        notifyCharacterIds: number[]
    ) {
        if (!notifyCharacterIds.length) {
            return;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        const routingKey = `fm.${worldId}.all.buddy`;
        const clientChannel = this.buddyChannelIndexForClient(channelIndex);
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, BUDDY_EVT_CHANNEL_UPDATE, {
            event_id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            occurred_at: new Date().toISOString(),
            character_id: characterId,
            channel_index: clientChannel,
            notify_character_ids: notifyCharacterIds,
        });
    }

    private async publishBuddyListUpdate(
        worldId: number,
        notifyCharacterIds: number[],
        syncAction: number,
        entries: BuddyListEntry[]
    ) {
        if (!notifyCharacterIds.length || !entries.length) {
            return;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        const routingKey = `fm.${worldId}.all.buddy`;
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, BUDDY_EVT_LIST_UPDATE, {
            event_id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            occurred_at: new Date().toISOString(),
            sync_action: syncAction,
            notify_character_ids: notifyCharacterIds,
            entries: entries.map((e) => ({
                character_id: e.characterId,
                name: e.name,
                group_name: e.groupName,
                pending: e.pending,
                channel_index: this.buddyChannelIndexForClient(e.channelIndex),
            })),
        });
    }

    private async publishBuddyAddRequest(
        worldId: number,
        targetCharacterId: number,
        fromCharacterId: number,
        fromName: string
    ) {
        if (!targetCharacterId || !fromCharacterId) {
            return;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        const routingKey = `fm.${worldId}.all.buddy`;
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, BUDDY_EVT_ADD_REQUEST, {
            event_id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            occurred_at: new Date().toISOString(),
            notify_character_ids: [targetCharacterId],
            from_character_id: fromCharacterId,
            from_name: fromName,
        });
    }

    private buddyChannelIndexForClient(channelIndex: number) {
        if (!Number.isInteger(channelIndex) || channelIndex < 0) {
            return -1;
        }
        return channelIndex;
    }

    async broadcastBuddyMultiChat(
        worldId: number,
        senderCharacterId: number,
        recipientCharacterIds: number[],
        senderName: string,
        message: string
    ): Promise<BroadcastBuddyMultiChatResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(senderCharacterId);
        const trimmedMsg = message;
        if (trimmedMsg.length <= 0 || trimmedMsg.length > MAX_MULTI_CHAT_MSG_LEN) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_UNKNOWN };
        }
        const name = senderName.trim();
        if (name.length <= 0 || name.length > MAX_MULTI_CHAT_NAME_LEN) {
            return { ok: false, code: BuddyErrorCode.BUDDY_ERROR_UNKNOWN };
        }

        const accepted = new Set(await this.buddyRepo.listAcceptedBuddyCharacterIds(worldId, senderCharacterId));
        let targets: number[];
        if (recipientCharacterIds.length > 0) {
            targets = recipientCharacterIds.filter((id) => id !== senderCharacterId && accepted.has(id));
        } else {
            targets = [...accepted].filter((id) => id !== senderCharacterId);
        }
        if (targets.length === 0) {
            return { ok: true, deliveredCount: 0 };
        }

        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        const routingKey = `fm.${worldId}.all.buddy`;
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, BUDDY_EVT_MULTI_CHAT, {
            event_id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            occurred_at: new Date().toISOString(),
            sender_character_id: senderCharacterId,
            chat_mode: 0,
            sender_name: name,
            message: trimmedMsg,
            notify_character_ids: targets,
        });
        return { ok: true, deliveredCount: targets.length };
    }

    private async buildEntry(worldId: number, buddy: CharacterBuddyModel): Promise<BuddyListEntry> {
        const buddyRow = await this.characterRepo.get(worldId, buddy.buddyCharacterId);
        const name = buddyRow?.name ?? "";
        let channelIndex = -1;
        if (!buddy.pending) {
            channelIndex = await this.sessionChannelIndex(worldId, buddy.buddyCharacterId);
        }
        return {
            characterId: buddy.buddyCharacterId,
            name,
            groupName: buddy.groupName,
            pending: buddy.pending,
            channelIndex,
        };
    }
}
