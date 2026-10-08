import {
    HiredMerchantItem,
    HiredMerchantSale,
    type HiredMerchant,
    type OpenHiredMerchantReply,
} from "../protobuf/generated/fminternal/internal_service";
import { HiredMerchantRepository, type HiredMerchantRow } from "../repos/hired-merchant-repository";
import { CharacterService, type SaveCharacterEntry } from "./character-service";

export class HiredMerchantService {
    private readonly merchantRepo: HiredMerchantRepository;
    private readonly characterService: CharacterService;

    constructor(hiredMerchantRepository: HiredMerchantRepository, characterService: CharacterService) {
        this.merchantRepo = hiredMerchantRepository;
        this.characterService = characterService;
    }

    private toProto(worldId: number, row: HiredMerchantRow): HiredMerchant {
        return {
            merchantId: row.merchant_id,
            worldId,
            accountId: row.account_id,
            characterId: row.character_id,
            ownerName: row.owner_name,
            channelId: row.channel_id ?? -1,
            mapId: row.map_id ?? 0,
            itemId: row.item_id,
            title: row.title,
            meso: row.meso,
            items: row.items.map((item) => HiredMerchantItem.fromJSON(item)),
            sold: row.sold.map((sale) => HiredMerchantSale.fromJSON(sale)),
            openedAtUnixMs: new Date(row.opened_at).getTime(),
            closedAtUnixMs: row.closed_at ? new Date(row.closed_at).getTime() : 0,
        };
    }

    async findHiredMerchant(worldId: number, accountId: number): Promise<HiredMerchant | undefined> {
        const row = await this.merchantRepo.findByAccount(worldId, accountId);
        if (row == null) {
            return undefined;
        }
        return this.toProto(worldId, row);
    }

    async openHiredMerchant(merchant: HiredMerchant): Promise<OpenHiredMerchantReply> {
        return this.merchantRepo.withAccountLock(merchant.worldId, merchant.accountId, async (txClient) => {
            const existing = await this.merchantRepo.findByAccount(merchant.worldId, merchant.accountId, { txClient });
            if (existing != null) {
                return { merchantId: 0, existing: this.toProto(merchant.worldId, existing) };
            }
            const closed = merchant.closedAtUnixMs !== 0;
            const row = await this.merchantRepo.insert(
                merchant.worldId,
                {
                    accountId: merchant.accountId,
                    characterId: merchant.characterId,
                    ownerName: merchant.ownerName,
                    channelId: closed ? null : merchant.channelId,
                    mapId: closed ? null : merchant.mapId,
                    itemId: merchant.itemId,
                    title: merchant.title,
                    meso: merchant.meso,
                    items: merchant.items.map((item) => HiredMerchantItem.toJSON(item)),
                    sold: merchant.sold.map((sale) => HiredMerchantSale.toJSON(sale)),
                    closedAt: closed ? new Date(merchant.closedAtUnixMs) : null,
                },
                { txClient }
            );
            return { merchantId: row.merchant_id, existing: undefined };
        });
    }

    async saveHiredMerchant(merchant: HiredMerchant, characters: SaveCharacterEntry[], close: boolean): Promise<void> {
        await this.merchantRepo.withAccountLock(merchant.worldId, merchant.accountId, async (txClient) => {
            if (characters.length > 0) {
                await this.characterService.saveCharacters(characters);
            }
            await this.merchantRepo.update(
                merchant.worldId,
                {
                    merchantId: merchant.merchantId,
                    accountId: merchant.accountId,
                    title: merchant.title,
                    meso: merchant.meso,
                    items: merchant.items.map((item) => HiredMerchantItem.toJSON(item)),
                    sold: merchant.sold.map((sale) => HiredMerchantSale.toJSON(sale)),
                },
                { txClient }
            );
            if (close === false) {
                return;
            }
            await this.merchantRepo.deleteEmpty(merchant.worldId, merchant.merchantId, merchant.accountId, { txClient });
            await this.merchantRepo.close(merchant.worldId, merchant.merchantId, merchant.accountId, { txClient });
        });
    }

    async claimStoreBank(worldId: number, accountId: number, characterId: number): Promise<HiredMerchant | undefined> {
        return this.merchantRepo.withAccountLock(worldId, accountId, async (txClient) => {
            const row = await this.merchantRepo.findByAccount(worldId, accountId, { txClient });
            if (row == null || row.closed_at == null || row.character_id !== characterId) {
                return undefined;
            }
            await this.merchantRepo.delete(worldId, row.merchant_id, accountId, { txClient });
            return this.toProto(worldId, row);
        });
    }

    async closeChannelMerchants(worldId: number, channelId: number): Promise<number> {
        return this.merchantRepo.closeChannel(worldId, channelId);
    }
}
