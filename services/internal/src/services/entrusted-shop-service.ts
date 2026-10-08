import {
    EntrustedShopItem,
    EntrustedShopSale,
    type EntrustedShop,
    type OpenEntrustedShopReply,
} from "../protobuf/generated/fminternal/internal_service";
import { EntrustedShopRepository, type EntrustedShopRow } from "../repos/entrusted-shop-repository";
import { CharacterService, type SaveCharacterEntry } from "./character-service";

export class EntrustedShopService {
    private readonly shopRepo: EntrustedShopRepository;
    private readonly characterService: CharacterService;

    constructor(entrustedShopRepository: EntrustedShopRepository, characterService: CharacterService) {
        this.shopRepo = entrustedShopRepository;
        this.characterService = characterService;
    }

    private toProto(worldId: number, row: EntrustedShopRow): EntrustedShop {
        return {
            shopId: row.shop_id,
            worldId,
            accountId: row.account_id,
            characterId: row.character_id,
            ownerName: row.owner_name,
            channelId: row.channel_id ?? -1,
            mapId: row.map_id ?? 0,
            itemId: row.item_id,
            title: row.title,
            meso: row.meso,
            items: row.items.map((item) => EntrustedShopItem.fromJSON(item)),
            sold: row.sold.map((sale) => EntrustedShopSale.fromJSON(sale)),
            openedAtUnixMs: new Date(row.opened_at).getTime(),
            closedAtUnixMs: row.closed_at ? new Date(row.closed_at).getTime() : 0,
        };
    }

    async findEntrustedShop(worldId: number, accountId: number): Promise<EntrustedShop | undefined> {
        const row = await this.shopRepo.findByAccount(worldId, accountId);
        if (row == null) {
            return undefined;
        }
        return this.toProto(worldId, row);
    }

    async openEntrustedShop(shop: EntrustedShop): Promise<OpenEntrustedShopReply> {
        return this.shopRepo.withAccountLock(shop.worldId, shop.accountId, async (txClient) => {
            const existing = await this.shopRepo.findByAccount(shop.worldId, shop.accountId, { txClient });
            if (existing != null) {
                return { shopId: 0, existing: this.toProto(shop.worldId, existing) };
            }
            const closed = shop.closedAtUnixMs !== 0;
            const row = await this.shopRepo.insert(
                shop.worldId,
                {
                    accountId: shop.accountId,
                    characterId: shop.characterId,
                    ownerName: shop.ownerName,
                    channelId: closed ? null : shop.channelId,
                    mapId: closed ? null : shop.mapId,
                    itemId: shop.itemId,
                    title: shop.title,
                    meso: shop.meso,
                    items: shop.items.map((item) => EntrustedShopItem.toJSON(item)),
                    sold: shop.sold.map((sale) => EntrustedShopSale.toJSON(sale)),
                    closedAt: closed ? new Date(shop.closedAtUnixMs) : null,
                },
                { txClient }
            );
            return { shopId: row.shop_id, existing: undefined };
        });
    }

    async saveEntrustedShop(shop: EntrustedShop, characters: SaveCharacterEntry[], close: boolean): Promise<void> {
        await this.shopRepo.withAccountLock(shop.worldId, shop.accountId, async (txClient) => {
            if (characters.length > 0) {
                await this.characterService.saveCharacters(characters);
            }
            await this.shopRepo.update(
                shop.worldId,
                {
                    shopId: shop.shopId,
                    accountId: shop.accountId,
                    title: shop.title,
                    meso: shop.meso,
                    items: shop.items.map((item) => EntrustedShopItem.toJSON(item)),
                    sold: shop.sold.map((sale) => EntrustedShopSale.toJSON(sale)),
                },
                { txClient }
            );
            if (close === false) {
                return;
            }
            await this.shopRepo.deleteEmpty(shop.worldId, shop.shopId, shop.accountId, { txClient });
            await this.shopRepo.close(shop.worldId, shop.shopId, shop.accountId, { txClient });
        });
    }

    async claimStoreBank(worldId: number, accountId: number, characterId: number): Promise<EntrustedShop | undefined> {
        return this.shopRepo.withAccountLock(worldId, accountId, async (txClient) => {
            const row = await this.shopRepo.findByAccount(worldId, accountId, { txClient });
            if (row == null || row.closed_at == null || row.character_id !== characterId) {
                return undefined;
            }
            await this.shopRepo.delete(worldId, row.shop_id, accountId, { txClient });
            return this.toProto(worldId, row);
        });
    }

    async closeChannelEntrustedShops(worldId: number, channelId: number): Promise<number> {
        return this.shopRepo.closeChannel(worldId, channelId);
    }
}
