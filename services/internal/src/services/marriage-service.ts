import type { PoolClient } from "pg";
import {
    Inventory,
    MarriageResult,
    MarriageStatus,
    ParcelResult,
    type ClaimWeddingGiftReply,
    type InviteWeddingGuestReply,
    type LoadWeddingGiftsReply,
    type Marriage,
    type MarriageReply,
    type Parcel,
    WeddingGift,
} from "../protobuf/generated/fminternal/internal_service";
import {
    MarriageRepository,
    MarriageRequestRepository,
    WeddingGiftClaimRequestRepository,
    WeddingGiftSendRequestRepository,
    type MarriageRow,
    type WeddingGiftRow,
} from "../repos/marriage-repository";
import { ParcelSendRequestRepository } from "../repos/parcel-repository";
import { RequestStatus } from "../repos/request-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { SessionRepository } from "../repos/session-repository";
import { ParcelService } from "./parcel-service";
import { RabbitMQService } from "./rabbitmq-service";

const DIVORCE_DELAY_MS = 72 * 60 * 60 * 1000;
const ENGAGEMENT_RING_BASE = 4210000;
const WEDDING_RING_BASE = 1112300;

export type GiveWeddingGiftInput = {
    worldId: number;
    receiverId: number;
    gift: WeddingGift;
    requestId: string;
};

export class MarriageService {
    private readonly marriageRepo: MarriageRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;
    private readonly marriageRequestRepo: MarriageRequestRepository;
    private readonly weddingGiftSendRequestRepo: WeddingGiftSendRequestRepository;
    private readonly weddingGiftClaimRequestRepo: WeddingGiftClaimRequestRepository;
    private readonly parcelSendRequestRepo: ParcelSendRequestRepository;
    private readonly parcelService: ParcelService;

    constructor(
        marriageRepository: MarriageRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService,
        marriageRequestRepository: MarriageRequestRepository,
        weddingGiftSendRequestRepository: WeddingGiftSendRequestRepository,
        weddingGiftClaimRequestRepository: WeddingGiftClaimRequestRepository,
        parcelSendRequestRepository: ParcelSendRequestRepository,
        parcelService: ParcelService
    ) {
        this.marriageRepo = marriageRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
        this.marriageRequestRepo = marriageRequestRepository;
        this.weddingGiftSendRequestRepo = weddingGiftSendRequestRepository;
        this.weddingGiftClaimRequestRepo = weddingGiftClaimRequestRepository;
        this.parcelSendRequestRepo = parcelSendRequestRepository;
        this.parcelService = parcelService;
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
            await this.rabbitmqService.publish(`fm.${worldId}.${channelId}.parcel`, "marriage_changed", {
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
        ringItemId: number,
        requestId: string
    ): Promise<MarriageReply> {
        return this.marriageRepo.withTransaction(worldId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock($1)", [0x6d617272]);
            if (requestId !== "0" && (await this.marriageRequestRepo.find(worldId, 0, requestId, { txClient: client })) != null) {
                return { result: MarriageResult.MARRIAGE_RESULT_ABORTED, marriage: undefined };
            }
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
            if (requestId !== "0") {
                await this.marriageRequestRepo.insert(worldId, 0, requestId, {}, client);
            }
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

    async reserveWedding(
        worldId: number,
        marriageId: number,
        characterId: number,
        ticketItemId: number,
        requestId: string
    ): Promise<MarriageReply> {
        const row = await this.marriageRepo.withTransaction(worldId, async (client) => {
            if (requestId !== "0" && (await this.marriageRequestRepo.find(worldId, 0, requestId, { txClient: client })) != null) {
                return MarriageResult.MARRIAGE_RESULT_ABORTED;
            }
            const found = await this.marriageRepo.find(client, marriageId);
            if (found == null || found.status !== MarriageStatus.MARRIAGE_STATUS_ENGAGED || found.ticket_item_id !== 0) {
                return MarriageResult.MARRIAGE_RESULT_INVALID_STATE;
            }
            if (found.groom_id !== characterId && found.bride_id !== characterId) {
                return MarriageResult.MARRIAGE_RESULT_INVALID_STATE;
            }
            found.ticket_item_id = ticketItemId;
            await this.marriageRepo.update(client, found);
            if (requestId !== "0") {
                await this.marriageRequestRepo.insert(worldId, 0, requestId, {}, client);
            }
            return found;
        });
        if (typeof row === "number") {
            return { result: row, marriage: undefined };
        }

        const partnerId = row.groom_id === characterId ? row.bride_id : row.groom_id;
        await this.notify(worldId, [partnerId], "booked");
        return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: this.toProto(row) };
    }

    async abortMarriageChange(worldId: number, requestId: string): Promise<boolean> {
        return this.marriageRepo.withTransaction(worldId, async (client) => {
            const request = await this.marriageRequestRepo.abort(worldId, 0, requestId, client);
            return request?.status === RequestStatus.Applied;
        });
    }

    async setWeddingWishlist(worldId: number, marriageId: number, characterId: number, wishes: string[], requestId: string): Promise<MarriageReply> {
        const row = await this.marriageRepo.withTransaction(worldId, async (client) => {
            const found = await this.marriageRepo.find(client, marriageId);
            if (found == null || found.ticket_item_id === 0) {
                return null;
            }
            const request = await this.marriageRequestRepo.find(worldId, 0, requestId, { txClient: client });
            if (request != null) {
                if (request.status === RequestStatus.Aborted || (found.groom_id !== characterId && found.bride_id !== characterId)) {
                    return null;
                }
                return found;
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
            await this.marriageRequestRepo.insert(worldId, 0, requestId, {}, client);
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

    // The marriage row lock serializes invites to one wedding. The invitation parcel commits on the guest's data shard
    // before the guest and the request id commit here, so a parcel recorded under the request id means the invite went through.
    async inviteWeddingGuest(worldId: number, marriageId: number, guestName: string, parcel: Parcel, requestId: string): Promise<InviteWeddingGuestReply> {
        return this.marriageRepo.withTransaction(worldId, async (client) => {
            const row = await this.marriageRepo.find(client, marriageId);
            const request = await this.marriageRequestRepo.find(worldId, 0, requestId, { txClient: client });
            if (request?.status === RequestStatus.Aborted) {
                return { result: MarriageResult.MARRIAGE_RESULT_ABORTED };
            }
            if (request != null) {
                return { result: MarriageResult.MARRIAGE_RESULT_OK };
            }

            const fail = async (result: MarriageResult): Promise<InviteWeddingGuestReply> => {
                await this.marriageRequestRepo.abort(worldId, 0, requestId, client);
                return { result };
            };
            const entry = await this.unifiedRepo.findCharacterNameEntry(guestName);
            if (entry == null || entry.world_id !== worldId) {
                return fail(MarriageResult.MARRIAGE_RESULT_GUEST_NOT_FOUND);
            }

            const sent = await this.parcelSendRequestRepo.find(worldId, entry.character_id, requestId);
            if (sent == null) {
                if (row == null || row.ticket_item_id === 0) {
                    return fail(MarriageResult.MARRIAGE_RESULT_INVALID_STATE);
                }
                if (entry.character_id === row.groom_id || entry.character_id === row.bride_id || row.guests.includes(entry.character_id)) {
                    return fail(MarriageResult.MARRIAGE_RESULT_GUEST_ALREADY_INVITED);
                }
                const result = await this.parcelService.sendParcel({
                    worldId,
                    recipientName: "",
                    parcel: { ...parcel, receiverId: entry.character_id },
                    senderAccountId: 0,
                    oneOfAKind: false,
                    requestId,
                    system: true,
                });
                if (result !== ParcelResult.PARCEL_RESULT_OK) {
                    return fail(MarriageResult.MARRIAGE_RESULT_INVALID_STATE);
                }
            }

            if (row != null && row.ticket_item_id !== 0 && row.guests.includes(entry.character_id) === false) {
                row.guests = [...row.guests, entry.character_id];
                await this.marriageRepo.update(client, row);
            }
            await this.marriageRequestRepo.insert(worldId, 0, requestId, {}, client);
            return { result: MarriageResult.MARRIAGE_RESULT_OK };
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
                return { result: MarriageResult.MARRIAGE_RESULT_OK, row, divorced: false, requested: true };
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
        if (loaded.requested) {
            const partnerId = loaded.row.groom_id === characterId ? loaded.row.bride_id : loaded.row.groom_id;
            await this.notify(worldId, [partnerId], "divorce_requested");
        }
        return { result: loaded.result, marriage: this.toProto(loaded.row) };
    }

    async cancelDivorce(worldId: number, marriageId: number, characterId: number): Promise<MarriageReply> {
        const row = await this.marriageRepo.withTransaction(worldId, async (client) => {
            const found = await this.marriageRepo.find(client, marriageId);
            if (found == null || found.status !== MarriageStatus.MARRIAGE_STATUS_MARRIED || found.divorce_requested_at == null) {
                return null;
            }
            if (found.groom_id !== characterId && found.bride_id !== characterId) {
                return null;
            }
            found.divorce_requester_id = 0;
            found.divorce_requested_at = null;
            await this.marriageRepo.update(client, found);
            return found;
        });
        if (row == null) {
            return { result: MarriageResult.MARRIAGE_RESULT_INVALID_STATE, marriage: undefined };
        }
        const partnerId = row.groom_id === characterId ? row.bride_id : row.groom_id;
        await this.notify(worldId, [partnerId], "divorce_canceled");
        return { result: MarriageResult.MARRIAGE_RESULT_OK, marriage: this.toProto(row) };
    }

    async giveWeddingGift(input: GiveWeddingGiftInput): Promise<MarriageResult> {
        return this.marriageRepo.withReceiverLock(input.worldId, input.receiverId, async (client) => {
            if (input.requestId !== "0") {
                const request = await this.weddingGiftSendRequestRepo.find(input.worldId, input.receiverId, input.requestId, { txClient: client });
                if (request?.status === RequestStatus.Aborted) {
                    return MarriageResult.MARRIAGE_RESULT_ABORTED;
                }
                if (request != null) {
                    return MarriageResult.MARRIAGE_RESULT_OK;
                }
            }

            await this.marriageRepo.insertGift(
                input.worldId,
                input.receiverId,
                input.gift.senderName,
                Inventory.toJSON(input.gift.item ?? Inventory.fromPartial({})),
                { txClient: client }
            );
            if (input.requestId !== "0") {
                await this.weddingGiftSendRequestRepo.insert(input.worldId, input.receiverId, input.requestId, {}, client);
            }
            return MarriageResult.MARRIAGE_RESULT_OK;
        });
    }

    async abortWeddingGift(worldId: number, receiverId: number, requestId: string): Promise<boolean> {
        return this.marriageRepo.withReceiverLock(worldId, receiverId, async (client) => {
            const request = await this.weddingGiftSendRequestRepo.abort(worldId, receiverId, requestId, client);
            return request?.status === RequestStatus.Applied;
        });
    }

    async loadWeddingGifts(worldId: number, characterId: number): Promise<LoadWeddingGiftsReply> {
        const rows = await this.marriageRepo.findGifts(worldId, characterId);
        return { gifts: rows.map((row) => this.giftToProto(row)) };
    }

    async claimWeddingGift(worldId: number, characterId: number, giftId: number, requestId: string): Promise<ClaimWeddingGiftReply> {
        return this.marriageRepo.withReceiverLock(worldId, characterId, async (client) => {
            if (requestId !== "0") {
                const claimed = await this.weddingGiftClaimRequestRepo.find(worldId, characterId, requestId, { txClient: client });
                if (claimed?.status === RequestStatus.Applied) {
                    return { result: MarriageResult.MARRIAGE_RESULT_OK, gift: WeddingGift.fromJSON(claimed.outcome.gift) };
                }
            }

            const row = await this.marriageRepo.deleteGift(client, characterId, giftId);
            if (row == null) {
                return { result: MarriageResult.MARRIAGE_RESULT_GIFT_NOT_FOUND, gift: undefined };
            }
            const gift = this.giftToProto(row);
            if (requestId !== "0") {
                await this.weddingGiftClaimRequestRepo.insert(worldId, characterId, requestId, { gift: WeddingGift.toJSON(gift) }, client);
            }
            return { result: MarriageResult.MARRIAGE_RESULT_OK, gift };
        });
    }

    async notifySpouseMap(worldId: number, characterId: number, spouseId: number, mapId: number, reply: boolean): Promise<void> {
        const channelId = await this.sessionRepo.findChannel(worldId, spouseId);
        if (channelId == null) {
            return;
        }
        await this.rabbitmqService.publish(`fm.${worldId}.${channelId}.parcel`, "spouse_moved", {
            character_id: spouseId,
            spouse_id: characterId,
            map_id: mapId,
            reply,
        });
    }
}
