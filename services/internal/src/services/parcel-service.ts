import {
    Inventory,
    ParcelResult,
    type CheckParcelArrivalsReply,
    type ClaimParcelReply,
    type LoadParcelsReply,
    Parcel,
} from "../protobuf/generated/fminternal/internal_service";
import { AppConfiguration } from "../config/app-configuration";
import { ParcelClaimRequestRepository, ParcelRepository, ParcelReturnRequestRepository, ParcelSendRequestRepository, type ParcelRow } from "../repos/parcel-repository";
import { RequestStatus } from "../repos/request-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { SessionRepository } from "../repos/session-repository";
import { RabbitMQService } from "./rabbitmq-service";

const ARRIVAL_DELAY_MS = 12 * 60 * 60 * 1000;
const KEEP_MS = 30 * 24 * 60 * 60 * 1000;

export type SendParcelInput = {
    worldId: number;
    recipientName: string;
    parcel: Parcel;
    senderAccountId: number;
    oneOfAKind: boolean;
    requestId: string;
    system: boolean;
};

export class ParcelService {
    private readonly app: AppConfiguration;
    private readonly parcelRepo: ParcelRepository;
    private readonly parcelSendRequestRepo: ParcelSendRequestRepository;
    private readonly parcelReturnRequestRepo: ParcelReturnRequestRepository;
    private readonly parcelClaimRequestRepo: ParcelClaimRequestRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;

    constructor(
        appConfiguration: AppConfiguration,
        parcelRepository: ParcelRepository,
        parcelSendRequestRepository: ParcelSendRequestRepository,
        parcelReturnRequestRepository: ParcelReturnRequestRepository,
        parcelClaimRequestRepository: ParcelClaimRequestRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService
    ) {
        this.app = appConfiguration;
        this.parcelRepo = parcelRepository;
        this.parcelSendRequestRepo = parcelSendRequestRepository;
        this.parcelReturnRequestRepo = parcelReturnRequestRepository;
        this.parcelClaimRequestRepo = parcelClaimRequestRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
    }

    private arrivedAt(row: ParcelRow): number {
        const sentAt = new Date(row.sent_at).getTime();
        if (row.quick) {
            return sentAt;
        }
        return sentAt + ARRIVAL_DELAY_MS;
    }

    private expired(row: ParcelRow, now: number): boolean {
        return this.arrivedAt(row) + KEEP_MS <= now;
    }

    private toProto(row: ParcelRow): Parcel {
        return {
            parcelId: row.parcel_id,
            receiverId: row.receiver_id,
            senderName: row.sender_name,
            meso: row.meso,
            quick: row.quick,
            message: row.message,
            sentAtUnixMs: new Date(row.sent_at).getTime(),
            item: row.item ? Inventory.fromJSON(row.item) : undefined,
        };
    }

    async loadParcels(worldId: number, characterId: number): Promise<LoadParcelsReply> {
        const now = Date.now();
        const rows = await this.parcelRepo.findByReceiver(worldId, characterId);
        const expired = rows.filter((row) => this.expired(row, now));
        await this.parcelRepo.deleteMany(worldId, characterId, expired.map((row) => row.parcel_id));
        return {
            parcels: rows.filter((row) => this.expired(row, now) === false).map((row) => this.toProto(row)),
            expired: expired.map((row) => this.toProto(row)),
        };
    }

    async sendParcel(input: SendParcelInput): Promise<ParcelResult> {
        let receiverId = input.parcel.receiverId;
        if (input.recipientName !== "") {
            const entry = await this.unifiedRepo.findCharacterNameEntry(input.recipientName);
            if (entry == null || entry.world_id !== input.worldId) {
                return ParcelResult.PARCEL_RESULT_RECIPIENT_NOT_FOUND;
            }
            if (input.system === false && input.senderAccountId !== 0 && entry.account_id === input.senderAccountId) {
                return ParcelResult.PARCEL_RESULT_SAME_ACCOUNT;
            }
            receiverId = entry.character_id;
        }

        const itemId = input.parcel.item?.itemId ?? 0;
        const result = await this.parcelRepo.withReceiverLock(input.worldId, receiverId, async (txClient) => {
            if (input.requestId !== "0") {
                const request = await this.parcelSendRequestRepo.find(input.worldId, receiverId, input.requestId, { txClient });
                if (request?.status === RequestStatus.Aborted) {
                    return ParcelResult.PARCEL_RESULT_ABORTED;
                }
                if (request != null) {
                    return ParcelResult.PARCEL_RESULT_OK;
                }
            }

            const now = Date.now();
            const rows = (await this.parcelRepo.findByReceiver(input.worldId, receiverId, { txClient })).filter(
                (row) => this.expired(row, now) === false
            );
            if (input.system === false && rows.length >= this.app.getMaxParcels()) {
                return ParcelResult.PARCEL_RESULT_RECIPIENT_FULL;
            }
            if (input.system === false && input.oneOfAKind && rows.some((row) => this.toProto(row).item?.itemId === itemId)) {
                return ParcelResult.PARCEL_RESULT_ONE_OF_A_KIND;
            }

            await this.parcelRepo.insert(
                input.worldId,
                {
                    receiverId,
                    senderName: input.parcel.senderName,
                    meso: input.parcel.meso,
                    quick: input.parcel.quick,
                    message: input.parcel.message,
                    item: input.parcel.item ? Inventory.toJSON(input.parcel.item) : null,
                },
                { txClient }
            );
            if (input.requestId !== "0") {
                await this.parcelSendRequestRepo.insert(input.worldId, receiverId, input.requestId, {}, txClient);
            }
            return ParcelResult.PARCEL_RESULT_OK;
        });
        if (result !== ParcelResult.PARCEL_RESULT_OK || input.parcel.quick === false) {
            return result;
        }

        // The parcel is already committed; a failed notice must not report the send as failed.
        try {
            const channelId = await this.sessionRepo.findChannel(input.worldId, receiverId);
            if (channelId == null || channelId < 0) {
                return result;
            }
            await this.rabbitmqService.publish(`fm.${input.worldId}.${channelId}.parcel`, "parcel_arrived", {
                character_id: receiverId,
                sender_name: input.parcel.senderName,
                quick: true,
            });
        } catch (err) {
            console.error("[parcel] arrival notice failed", err);
        }
        return result;
    }

    async returnParcels(worldId: number, characterId: number, requestId: string, parcels: Parcel[]): Promise<ParcelResult> {
        let created = false;
        const stored = await this.parcelRepo.withReceiverLock(worldId, characterId, async (txClient) => {
            if (requestId !== "0") {
                const request = await this.parcelReturnRequestRepo.find(worldId, characterId, requestId, { txClient });
                if (request != null) {
                    return ParcelResult.PARCEL_RESULT_OK;
                }
            }

            const inserted: Parcel[] = [];
            for (const parcel of parcels) {
                if (parcel.item == null && parcel.meso === 0) {
                    continue;
                }
                const row = await this.parcelRepo.insert(
                    worldId,
                    {
                        receiverId: characterId,
                        senderName: parcel.senderName,
                        meso: parcel.meso,
                        quick: true,
                        message: parcel.message,
                        item: parcel.item ? Inventory.toJSON(parcel.item) : null,
                    },
                    { txClient }
                );
                inserted.push(this.toProto(row));
            }
            if (inserted.length === 0) {
                return ParcelResult.PARCEL_RESULT_OK;
            }
            if (requestId !== "0") {
                await this.parcelReturnRequestRepo.insert(worldId, characterId, requestId, {}, txClient);
            }
            created = true;
            return ParcelResult.PARCEL_RESULT_OK;
        });
        if (stored !== ParcelResult.PARCEL_RESULT_OK || created === false) {
            return stored;
        }

        try {
            const channelId = await this.sessionRepo.findChannel(worldId, characterId);
            if (channelId == null || channelId < 0) {
                return stored;
            }
            await this.rabbitmqService.publish(`fm.${worldId}.${channelId}.parcel`, "parcel_arrived", {
                character_id: characterId,
                sender_name: parcels[0]?.senderName ?? "",
                quick: true,
            });
        } catch (err) {
            console.error("[parcel] return notice failed", err);
        }
        return stored;
    }

    async abortParcelSend(worldId: number, recipientName: string, requestId: string): Promise<boolean> {
        const entry = await this.unifiedRepo.findCharacterNameEntry(recipientName);
        if (entry == null || entry.world_id !== worldId) {
            return false;
        }

        return this.parcelRepo.withReceiverLock(worldId, entry.character_id, async (txClient) => {
            const request = await this.parcelSendRequestRepo.abort(worldId, entry.character_id, requestId, txClient);
            return request?.status === RequestStatus.Applied;
        });
    }

    async claimParcel(worldId: number, characterId: number, parcelId: number, requestId: string): Promise<ClaimParcelReply> {
        return this.parcelRepo.withReceiverLock(worldId, characterId, async (txClient) => {
            if (requestId !== "0") {
                const claimed = await this.parcelClaimRequestRepo.find(worldId, characterId, requestId, { txClient });
                if (claimed?.status === RequestStatus.Applied) {
                    return { result: ParcelResult.PARCEL_RESULT_OK, parcel: Parcel.fromJSON(claimed.outcome.parcel) };
                }
            }

            const now = Date.now();
            const row = (await this.parcelRepo.findByReceiver(worldId, characterId, { txClient })).find(
                (x) => x.parcel_id === parcelId
            );
            if (row == null || this.expired(row, now)) {
                return { result: ParcelResult.PARCEL_RESULT_NOT_FOUND, parcel: undefined };
            }
            if (this.arrivedAt(row) > now) {
                return { result: ParcelResult.PARCEL_RESULT_NOT_ARRIVED, parcel: undefined };
            }
            await this.parcelRepo.delete(worldId, characterId, parcelId, { txClient });
            const parcel = this.toProto(row);
            if (requestId !== "0") {
                await this.parcelClaimRequestRepo.insert(worldId, characterId, requestId, { parcel: Parcel.toJSON(parcel) }, txClient);
            }
            return { result: ParcelResult.PARCEL_RESULT_OK, parcel };
        });
    }

    async deleteParcel(worldId: number, characterId: number, parcelId: number): Promise<ParcelResult> {
        const row = await this.parcelRepo.delete(worldId, characterId, parcelId);
        if (row == null) {
            return ParcelResult.PARCEL_RESULT_NOT_FOUND;
        }
        return ParcelResult.PARCEL_RESULT_OK;
    }

    async checkParcelArrivals(worldId: number, characterId: number): Promise<CheckParcelArrivalsReply> {
        const now = Date.now();
        const rows = (await this.parcelRepo.findByReceiver(worldId, characterId)).filter(
            (row) => row.notified === false && this.arrivedAt(row) <= now && this.expired(row, now) === false
        );
        const first = rows[0];
        if (first === undefined) {
            return { count: 0, senderName: "", quick: false };
        }
        await this.parcelRepo.markNotified(worldId, characterId, rows.map((row) => row.parcel_id));
        return { count: rows.length, senderName: first.sender_name, quick: first.quick };
    }
}
