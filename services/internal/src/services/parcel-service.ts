import {
    InventoryPersisted,
    ParcelResult,
    type CheckParcelArrivalsReply,
    type ClaimParcelReply,
    type LoadParcelsReply,
    type ParcelPersisted,
} from "../protobuf/generated/fminternal/internal_service";
import { AppConfiguration } from "../config/app-configuration";
import { ParcelRepository, type ParcelRow } from "../repos/parcel-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { SessionRepository } from "../repos/session-repository";
import { RabbitMQService } from "./rabbitmq-service";
import { CharacterService, type SaveCharacterEntry } from "./character-service";

const DELIVERY_DELAY_MS = 12 * 60 * 60 * 1000;
const KEEP_MS = 30 * 24 * 60 * 60 * 1000;
const AMQ_DIRECT_EXCHANGE = "amq.direct";

export type SendParcelInput = {
    worldId: number;
    recipientName: string;
    parcel: ParcelPersisted;
    senderAccountId: number;
    sender?: SaveCharacterEntry;
    oneOfAKind: boolean;
};

export class ParcelService {
    private readonly app: AppConfiguration;
    private readonly parcelRepo: ParcelRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;
    private readonly characterService: CharacterService;

    constructor(
        appConfiguration: AppConfiguration,
        parcelRepository: ParcelRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService,
        characterService: CharacterService
    ) {
        this.app = appConfiguration;
        this.parcelRepo = parcelRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
        this.characterService = characterService;
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

    private toProto(row: ParcelRow): ParcelPersisted {
        return {
            parcelId: row.parcel_id,
            receiverId: row.receiver_id,
            senderName: row.sender_name,
            meso: row.meso,
            quick: row.quick,
            message: row.message,
            sentAtUnixMs: new Date(row.sent_at).getTime(),
            item: row.item ? InventoryPersisted.fromJSON(row.item) : undefined,
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

            if (input.sender) {
                await this.characterService.saveCharacters([input.sender]);
            }
            await this.parcelRepo.insert(
                input.worldId,
                {
                    receiverId,
                    senderName: input.parcel.senderName,
                    meso: input.parcel.meso,
                    quick: input.parcel.quick,
                    message: input.parcel.message,
                    item: input.parcel.item ? InventoryPersisted.toJSON(input.parcel.item) : null,
                },
                { txClient }
            );
            return ParcelResult.PARCEL_RESULT_OK;
        });
        if (result !== ParcelResult.PARCEL_RESULT_OK || input.parcel.quick === false) {
            return result;
        }

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
        return result;
    }

    async claimParcel(worldId: number, characterId: number, parcelId: number): Promise<ClaimParcelReply> {
        return this.parcelRepo.withReceiverLock(worldId, characterId, async (txClient) => {
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
            return { result: ParcelResult.PARCEL_RESULT_OK, parcel: this.toProto(row) };
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
