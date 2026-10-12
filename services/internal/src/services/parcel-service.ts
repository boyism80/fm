import {
    Inventory,
    ParcelResult,
    type CheckParcelArrivalsReply,
    type ClaimParcelReply,
    type LoadParcelsReply,
    Parcel,
} from "../protobuf/generated/fminternal/internal_service";
import { AppConfiguration } from "../config/app-configuration";
import { ParcelRepository, type ParcelRow } from "../repos/parcel-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { SessionRepository } from "../repos/session-repository";
import { RabbitMQService } from "./rabbitmq-service";

const DELIVERY_DELAY_MS = 12 * 60 * 60 * 1000;
const KEEP_MS = 30 * 24 * 60 * 60 * 1000;
const AMQ_DIRECT_EXCHANGE = "amq.direct";

export type SendParcelInput = {
    worldId: number;
    recipientName: string;
    parcel: Parcel;
    senderAccountId: number;
    oneOfAKind: boolean;
    requestId: string;
};

export class ParcelService {
    private readonly app: AppConfiguration;
    private readonly parcelRepo: ParcelRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;

    constructor(
        appConfiguration: AppConfiguration,
        parcelRepository: ParcelRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService
    ) {
        this.app = appConfiguration;
        this.parcelRepo = parcelRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
    }

    private arrivedAt(row: ParcelRow): number {
        const sentAt = new Date(row.sent_at).getTime();
        if (row.quick) {
            return sentAt;
        }
        return sentAt + DELIVERY_DELAY_MS;
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
            if (input.senderAccountId !== 0 && entry.account_id === input.senderAccountId) {
                return ParcelResult.PARCEL_RESULT_SAME_ACCOUNT;
            }
            receiverId = entry.character_id;
        }

        const itemId = input.parcel.item?.itemId ?? 0;
        const result = await this.parcelRepo.withReceiverLock(input.worldId, receiverId, async (txClient) => {
            if (input.requestId !== "0") {
                const request = await this.parcelRepo.findSendRequest(input.worldId, receiverId, input.requestId, { txClient });
                if (request?.aborted) {
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
            if (rows.length >= this.app.getMaxParcels()) {
                return ParcelResult.PARCEL_RESULT_RECIPIENT_FULL;
            }
            if (input.oneOfAKind && rows.some((row) => this.toProto(row).item?.itemId === itemId)) {
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
                await this.parcelRepo.insertSendRequest(input.worldId, receiverId, input.requestId, false, { txClient });
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
            await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
            await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, `fm.${input.worldId}.${channelId}.parcel`, "parcel_arrived", {
                occurred_at: new Date().toISOString(),
                character_id: receiverId,
                sender_name: input.parcel.senderName,
                quick: true,
            });
        } catch (err) {
            console.error("[parcel] arrival notice failed", err);
        }
        return result;
    }

    async abortParcelSend(worldId: number, recipientName: string, requestId: string): Promise<boolean> {
        const entry = await this.unifiedRepo.findCharacterNameEntry(recipientName);
        if (entry == null || entry.world_id !== worldId) {
            return false;
        }

        return this.parcelRepo.withReceiverLock(worldId, entry.character_id, async (txClient) => {
            const request = await this.parcelRepo.findSendRequest(worldId, entry.character_id, requestId, { txClient });
            if (request != null) {
                return request.aborted === false;
            }

            await this.parcelRepo.insertSendRequest(worldId, entry.character_id, requestId, true, { txClient });
            return false;
        });
    }

    async claimParcel(worldId: number, characterId: number, parcelId: number, requestId: string): Promise<ClaimParcelReply> {
        return this.parcelRepo.withReceiverLock(worldId, characterId, async (txClient) => {
            if (requestId !== "0") {
                const claimed = await this.parcelRepo.findClaimRequest(worldId, characterId, requestId, { txClient });
                if (claimed != null) {
                    return { result: ParcelResult.PARCEL_RESULT_OK, parcel: Parcel.fromJSON(claimed) };
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
                await this.parcelRepo.insertClaimRequest(worldId, characterId, requestId, Parcel.toJSON(parcel), { txClient });
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
