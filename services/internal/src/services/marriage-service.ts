import type { PoolClient } from "pg";
import {
    Inventory,
    MarriageResult,
    MarriageStatus,
    type ClaimWeddingGiftReply,
    type InviteWeddingGuestReply,
    type LoadWeddingGiftsReply,
    type Marriage,
    type MarriageReply,
    type WeddingGift,
} from "../protobuf/generated/fminternal/internal_service";
import { MarriageRepository, type MarriageRow, type WeddingGiftRow } from "../repos/marriage-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { SessionRepository } from "../repos/session-repository";
import { RabbitMQService } from "./rabbitmq-service";
import { CharacterService, type SaveCharacterEntry } from "./character-service";

const DIVORCE_DELAY_MS = 72 * 60 * 60 * 1000;
const ENGAGEMENT_RING_BASE = 4210000;
const WEDDING_RING_BASE = 1112300;
const AMQ_DIRECT_EXCHANGE = "amq.direct";

export type GiveWeddingGiftInput = {
    worldId: number;
    receiverId: number;
    gift: WeddingGift;
    sender?: SaveCharacterEntry;
};

export class MarriageService {
    private readonly marriageRepo: MarriageRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;
    private readonly characterService: CharacterService;

    constructor(
        marriageRepository: MarriageRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService,
        characterService: CharacterService
    ) {
        this.marriageRepo = marriageRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
        this.characterService = characterService;
    }

    private toProto(row: MarriageRow): Marriage {
        return {
            marriageId: row.marriage_id,
            groomId: row.groom_id,
            brideId: row.bride_id,
            groomName: row.groom_name,
            brideName: row.bride_name,
            groomItemId: row.groom_item_id,
            brideItemId: row.bride_item_id,
            status: row.status,
            ticketItemId: row.ticket_item_id,
            groomWished: row.groom_wished,
            brideWished: row.bride_wished,
            groomWishes: row.groom_wishes,
            brideWishes: row.bride_wishes,
            guests: row.guests,
            divorceRequesterId: row.divorce_requester_id,
            divorceRequestedAtUnixMs: row.divorce_requested_at ? new Date(row.divorce_requested_at).getTime() : 0,
        };
    }

    private giftToProto(row: WeddingGiftRow): WeddingGift {
        return {
            giftId: row.gift_id,
            senderName: row.sender_name,
            item: Inventory.fromJSON(row.item),
        };
    }

    private divorceDue(row: MarriageRow, nowMs: number): boolean {
        if (row.divorce_requested_at == null) {
            return false;
        }
        return new Date(row.divorce_requested_at).getTime() + DIVORCE_DELAY_MS <= nowMs;
    }

    private async notify(worldId: number, characterIds: number[], event: string): Promise<void> {
        for (const characterId of characterIds) {
            const channelId = await this.sessionRepo.findChannel(worldId, characterId);
            if (channelId == null) {
                continue;
            }
            await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
            await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, `fm.${worldId}.${channelId}.parcel`, "marriage_changed", {
                occurred_at: new Date().toISOString(),
                character_id: characterId,
                event,
            });
        }
    }

    private async loadActive(worldId: number, find: (client: PoolClient) => Promise<MarriageRow | null>): Promise<MarriageRow | null> {
        const loaded = await this.marriageRepo.withTransaction(worldId, async (client) => {
            const row = await find(client);
            if (row == null) {
                return { row: null, divorced: null };
            }
            if (this.divorceDue(row, Date.now())) {
                await this.marriageRepo.end(client, row.marriage_id);
                return { row: null, divorced: row };
            }
            return { row, divorced: null };
        });
        if (loaded.divorced != null) {
            await this.notify(worldId, [loaded.divorced.groom_id, loaded.divorced.bride_id], "divorced");
        }
        return loaded.row;
    }

    async loadMarriage(worldId: number, characterId: number): Promise<MarriageReply> {
        const row = await this.loadActive(worldId, (client) => this.marriageRepo.findByCharacter(client, characterId));
        if (row == null) {
            return { result: MarriageResult.MARRIAGE_RESULT_NOT_FOUND, marriage: undefined };
        }
        return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: this.toProto(row) };
    }

    async getMarriage(worldId: number, marriageId: number): Promise<MarriageReply> {
        const row = await this.loadActive(worldId, (client) => this.marriageRepo.find(client, marriageId));
        if (row == null) {
            return { result: MarriageResult.MARRIAGE_RESULT_NOT_FOUND, marriage: undefined };
        }
        return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: this.toProto(row) };
    }

    async createMarriage(
        worldId: number,
        groomId: number,
        groomName: string,
        brideId: number,
        brideName: string,
        ringItemId: number
    ): Promise<MarriageReply> {
        return this.marriageRepo.withTransaction(worldId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock($1)", [0x6d617272]);
            if ((await this.marriageRepo.findByCharacter(client, groomId)) != null) {
                return { result: MarriageResult.MARRIAGE_RESULT_ALREADY_ENGAGED, marriage: undefined };
            }
            if ((await this.marriageRepo.findByCharacter(client, brideId)) != null) {
                return { result: MarriageResult.MARRIAGE_RESULT_PARTNER_ENGAGED, marriage: undefined };
            }
            const row = await this.marriageRepo.insert(client, {
                groomId,
                brideId,
                groomName,
                brideName,
                groomItemId: ringItemId,
                brideItemId: ringItemId,
                status: MarriageStatus.MARRIAGE_STATUS_ENGAGED,
            });
            return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: this.toProto(row) };
        });
    }

    async breakEngagement(worldId: number, marriageId: number, characterId: number): Promise<MarriageReply> {
        const row = await this.marriageRepo.withTransaction(worldId, async (client) => {
            const found = await this.marriageRepo.find(client, marriageId);
            if (found == null || found.status !== MarriageStatus.MARRIAGE_STATUS_ENGAGED) {
                return null;
            }
            if (found.groom_id !== characterId && found.bride_id !== characterId) {
                return null;
            }
            await this.marriageRepo.end(client, marriageId);
            return found;
        });
        if (row == null) {
            return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, marriage: undefined };
        }

        const partnerId = row.groom_id === characterId ? row.bride_id : row.groom_id;
        await this.notify(worldId, [partnerId], "broken");
        return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: undefined };
    }

    async reserveWedding(worldId: number, marriageId: number, ticketItemId: number): Promise<MarriageReply> {
        return this.marriageRepo.withTransaction(worldId, async (client) => {
            const row = await this.marriageRepo.find(client, marriageId);
            if (row == null || row.status !== MarriageStatus.MARRIAGE_STATUS_ENGAGED || row.ticket_item_id !== 0) {
                return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, marriage: undefined };
            }
            row.ticket_item_id = ticketItemId;
            await this.marriageRepo.update(client, row);
            return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: this.toProto(row) };
        });
    }

    async setWeddingWishlist(worldId: number, marriageId: number, characterId: number, wishes: string[]): Promise<MarriageReply> {
        const row = await this.marriageRepo.withTransaction(worldId, async (client) => {
            const found = await this.marriageRepo.find(client, marriageId);
            if (found == null || found.ticket_item_id === 0) {
                return null;
            }
            if (found.groom_id === characterId && found.groom_wished === false) {
                found.groom_wished = true;
                found.groom_wishes = wishes;
            } else if (found.bride_id === characterId && found.bride_wished === false) {
                found.bride_wished = true;
                found.bride_wishes = wishes;
            } else {
                return null;
            }
            await this.marriageRepo.update(client, found);
            return found;
        });
        if (row == null) {
            return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, marriage: undefined };
        }

        if (row.groom_wished && row.bride_wished) {
            const partnerId = row.groom_id === characterId ? row.bride_id : row.groom_id;
            await this.notify(worldId, [partnerId], "reserved");
        }
        return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: this.toProto(row) };
    }

    async inviteWeddingGuest(worldId: number, marriageId: number, guestName: string): Promise<InviteWeddingGuestReply> {
        const entry = await this.unifiedRepo.findCharacterNameEntry(guestName);
        if (entry == null || entry.world_id !== worldId) {
            return { result: MarriageResult.MARRIAGE_RESULT_GUEST_NOT_FOUND, guestId: 0, guestName: "" };
        }

        return this.marriageRepo.withTransaction(worldId, async (client) => {
            const row = await this.marriageRepo.find(client, marriageId);
            if (row == null || row.ticket_item_id === 0) {
                return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, guestId: 0, guestName: "" };
            }
            if (entry.character_id === row.groom_id || entry.character_id === row.bride_id || row.guests.includes(entry.character_id)) {
                return { result: MarriageResult.MARRIAGE_RESULT_GUEST_ALREADY_INVITED, guestId: 0, guestName: "" };
            }
            row.guests = [...row.guests, entry.character_id];
            await this.marriageRepo.update(client, row);
            return { result: MarriageResult.MARRIAGE_RESULT_OK, guestId: entry.character_id, guestName: entry.name };
        });
    }

    async finishWedding(worldId: number, marriageId: number): Promise<MarriageReply> {
        const row = await this.marriageRepo.withTransaction(worldId, async (client) => {
            const found = await this.marriageRepo.find(client, marriageId);
            if (found == null || found.status !== MarriageStatus.MARRIAGE_STATUS_ENGAGED || found.ticket_item_id === 0) {
                return null;
            }
            found.status = MarriageStatus.MARRIAGE_STATUS_MARRIED;
            found.groom_item_id = WEDDING_RING_BASE + (found.groom_item_id - ENGAGEMENT_RING_BASE);
            found.bride_item_id = WEDDING_RING_BASE + (found.bride_item_id - ENGAGEMENT_RING_BASE);
            found.ticket_item_id = 0;
            found.groom_wished = false;
            found.bride_wished = false;
            found.groom_wishes = [];
            found.bride_wishes = [];
            found.guests = [];
            await this.marriageRepo.update(client, found);
            return found;
        });
        if (row == null) {
            return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, marriage: undefined };
        }

        await this.notify(worldId, [row.groom_id, row.bride_id], "married");
        return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: this.toProto(row) };
    }

    async requestDivorce(worldId: number, marriageId: number, characterId: number, nowMs: number): Promise<MarriageReply> {
        const loaded = await this.marriageRepo.withTransaction(worldId, async (client) => {
            const row = await this.marriageRepo.find(client, marriageId);
            if (row == null || row.status !== MarriageStatus.MARRIAGE_STATUS_MARRIED) {
                return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, row: null, divorced: false };
            }
            if (row.groom_id !== characterId && row.bride_id !== characterId) {
                return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, row: null, divorced: false };
            }
            if (row.divorce_requested_at == null) {
                row.divorce_requester_id = characterId;
                row.divorce_requested_at = new Date(nowMs);
                await this.marriageRepo.update(client, row);
                return { result: MarriageResult.MARRIAGE_RESULT_OK, row, divorced: false };
            }
            if (this.divorceDue(row, nowMs) === false) {
                return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, row, divorced: false };
            }
            await this.marriageRepo.end(client, marriageId);
            return { result: MarriageResult.MARRIAGE_RESULT_OK, row, divorced: true };
        });
        if (loaded.row == null) {
            return { result: loaded.result, marriage: undefined };
        }
        if (loaded.divorced) {
            await this.notify(worldId, [loaded.row.groom_id, loaded.row.bride_id], "divorced");
            return { result: loaded.result, marriage: undefined };
        }
        return { result: loaded.result, marriage: this.toProto(loaded.row) };
    }

    async giveWeddingGift(input: GiveWeddingGiftInput): Promise<MarriageResult> {
        if (input.sender) {
            await this.characterService.saveCharacters([input.sender]);
        }
        await this.marriageRepo.insertGift(
            input.worldId,
            input.receiverId,
            input.gift.senderName,
            Inventory.toJSON(input.gift.item ?? Inventory.fromPartial({}))
        );
        return MarriageResult.MARRIAGE_RESULT_OK;
    }

    async loadWeddingGifts(worldId: number, characterId: number): Promise<LoadWeddingGiftsReply> {
        const rows = await this.marriageRepo.findGifts(worldId, characterId);
        return { gifts: rows.map((row) => this.giftToProto(row)) };
    }

    async claimWeddingGift(worldId: number, characterId: number, giftId: number): Promise<ClaimWeddingGiftReply> {
        const row = await this.marriageRepo.deleteGift(worldId, characterId, giftId);
        if (row == null) {
            return { result: MarriageResult.MARRIAGE_RESULT_GIFT_NOT_FOUND, gift: undefined };
        }
        return { result: MarriageResult.MARRIAGE_RESULT_OK, gift: this.giftToProto(row) };
    }

    async notifySpouseMap(worldId: number, characterId: number, spouseId: number, mapId: number, reply: boolean): Promise<void> {
        const channelId = await this.sessionRepo.findChannel(worldId, spouseId);
        if (channelId == null) {
            return;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, `fm.${worldId}.${channelId}.parcel`, "spouse_moved", {
            occurred_at: new Date().toISOString(),
            character_id: spouseId,
            spouse_id: characterId,
            map_id: mapId,
            reply,
        });
    }
}
