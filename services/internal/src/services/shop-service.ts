import {
    ShopItem,
    ShopSale,
    type OpenShopReply,
    type Shop,
    type ShopSearchEntry,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import { ShopOpenRequestRepository, ShopRepository, type ShopRow } from "@fm/shared/repos/shop-repository";
import { RequestStatus } from "@fm/shared/repos/request-repository";
import { InternalContext } from "@fm/shared/context/internal-context";
import { CharacterService, type SaveCharacterEntry } from "./character-service";

const ShopKindEntrusted = 5;
const SearchLimit = 200;
const PopularLimit = 10;

// The search counts are the ranking itself, not a cache, so they live in the world's Redis without the cache prefix
// and outlast a restart of any internal process.
export class ShopSearchRanking {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    private key(worldId: number): string {
        return `w${worldId}:shop-search-count`;
    }

    async add(worldId: number, itemId: number): Promise<void> {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.zincrby(this.key(worldId), 1, String(itemId));
    }

    async top(worldId: number, n: number): Promise<number[]> {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        return (await client.zrevrange(this.key(worldId), 0, n - 1)).map(Number);
    }
}

export class ShopService {
    private readonly shopRepo: ShopRepository;
    private readonly shopOpenRequestRepo: ShopOpenRequestRepository;
    private readonly characterService: CharacterService;
    private readonly ranking: ShopSearchRanking;

    constructor(
        shopRepository: ShopRepository,
        shopOpenRequestRepository: ShopOpenRequestRepository,
        characterService: CharacterService,
        shopSearchRanking: ShopSearchRanking,
    ) {
        this.shopRepo = shopRepository;
        this.shopOpenRequestRepo = shopOpenRequestRepository;
        this.characterService = characterService;
        this.ranking = shopSearchRanking;
    }

    private toProto(worldId: number, row: ShopRow): Shop {
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
            items: row.items.map((item) => ShopItem.fromJSON(item)),
            sold: row.sold.map((sale) => ShopSale.fromJSON(sale)),
            openedAtUnixMs: new Date(row.opened_at).getTime(),
            closedAtUnixMs: row.closed_at ? new Date(row.closed_at).getTime() : 0,
            kind: row.kind,
            sn: row.sn ?? 0,
            storeBankId: row.store_bank_id ?? "0",
        };
    }

    async findEntrustedShop(
        worldId: number,
        accountId: number,
    ): Promise<{ shop: Shop | undefined; storeBank: Shop[] }> {
        const rows = await this.shopRepo.findByAccount(worldId, accountId);
        const open = rows.find((row) => row.kind === ShopKindEntrusted && row.closed_at == null);
        return {
            shop: open ? this.toProto(worldId, open) : undefined,
            storeBank: rows.filter((row) => row.closed_at != null).map((row) => this.toProto(worldId, row)),
        };
    }

    async openShop(shop: Shop, requestId: string): Promise<OpenShopReply> {
        if (requestId === "0") {
            throw Object.assign(new Error("request_id is required"), { code: "INVALID_PAYLOAD" });
        }
        return this.shopRepo.withAccountLock(shop.worldId, shop.accountId, async (txClient) => {
            const request = await this.shopOpenRequestRepo.find(shop.worldId, shop.accountId, requestId, { txClient });
            if (request?.status === RequestStatus.Aborted) {
                return { shopId: 0, existing: undefined, aborted: true };
            }
            if (request != null) {
                return { shopId: request.outcome.shop_id, existing: undefined, aborted: false };
            }
            const rows = await this.shopRepo.findByAccount(shop.worldId, shop.accountId, { txClient });
            const existing = rows.find((row) => {
                if (shop.kind === ShopKindEntrusted) {
                    return row.closed_at != null || row.kind === ShopKindEntrusted;
                }
                return row.closed_at == null && row.kind === shop.kind && row.character_id === shop.characterId;
            });
            if (existing != null) {
                return { shopId: 0, existing: this.toProto(shop.worldId, existing), aborted: false };
            }
            const row = await this.shopRepo.insert(
                shop.worldId,
                {
                    kind: shop.kind,
                    accountId: shop.accountId,
                    characterId: shop.characterId,
                    ownerName: shop.ownerName,
                    channelId: shop.channelId,
                    mapId: shop.mapId,
                    itemId: shop.itemId,
                    title: shop.title,
                    meso: shop.meso,
                    items: shop.items.map((item) => ShopItem.toJSON(item)),
                    sold: shop.sold.map((sale) => ShopSale.toJSON(sale)),
                    closedAt: null,
                    storeBankId: null,
                },
                { txClient },
            );
            if (!row) {
                throw new Error("shop insert returned no row");
            }
            await this.shopOpenRequestRepo.insert(
                shop.worldId,
                shop.accountId,
                requestId,
                { shop_id: Number(row.shop_id) },
                txClient,
            );
            return { shopId: row.shop_id, existing: undefined, aborted: false };
        });
    }

    async abortOpenShop(worldId: number, accountId: number, requestId: string): Promise<boolean> {
        if (requestId === "0") {
            return false;
        }
        return this.shopRepo.withAccountLock(worldId, accountId, async (txClient) => {
            const request = await this.shopOpenRequestRepo.abort(worldId, accountId, requestId, txClient);
            if (request == null) {
                return false;
            }
            if (request.status === RequestStatus.Aborted) {
                return true;
            }
            if ((await this.shopRepo.abortOpen(worldId, accountId, request.outcome.shop_id, { txClient })) === false) {
                return false;
            }
            await this.shopOpenRequestRepo.transition(
                worldId,
                accountId,
                requestId,
                RequestStatus.Applied,
                RequestStatus.Aborted,
                txClient,
            );
            return true;
        });
    }

    async saveShop(shop: Shop, characters: SaveCharacterEntry[], close: boolean, storeBank: Shop[]): Promise<void> {
        if (characters.length > 0) {
            await this.characterService.saveCharacters(characters);
        }
        await this.shopRepo.withAccountLock(shop.worldId, shop.accountId, async (txClient) => {
            for (const kept of storeBank) {
                await this.shopRepo.insert(
                    shop.worldId,
                    {
                        kind: kept.kind,
                        accountId: shop.accountId,
                        characterId: kept.characterId,
                        ownerName: kept.ownerName,
                        channelId: null,
                        mapId: null,
                        itemId: kept.itemId,
                        title: kept.title,
                        meso: kept.meso,
                        items: kept.items.map((item) => ShopItem.toJSON(item)),
                        sold: [],
                        closedAt: new Date(),
                        storeBankId: kept.storeBankId,
                    },
                    { txClient },
                );
            }
            await this.shopRepo.update(
                shop.worldId,
                {
                    shopId: shop.shopId,
                    accountId: shop.accountId,
                    sn: shop.sn === 0 ? null : shop.sn,
                    title: shop.title,
                    meso: shop.meso,
                    items: shop.items.map((item) => ShopItem.toJSON(item)),
                    sold: shop.sold.map((sale) => ShopSale.toJSON(sale)),
                },
                { txClient },
            );
            if (close === false) {
                return;
            }
            await this.shopRepo.deleteEmpty(shop.worldId, shop.shopId, shop.accountId, { txClient });
            await this.shopRepo.close(shop.worldId, shop.shopId, shop.accountId, { txClient });
        });
    }

    async closeChannelShops(worldId: number, channelId: number): Promise<number> {
        return this.shopRepo.closeChannel(worldId, channelId);
    }

    async searchShops(worldId: number, itemId: number, highFirst: boolean): Promise<ShopSearchEntry[]> {
        const entries: ShopSearchEntry[] = [];
        for (const row of await this.shopRepo.findOpen(worldId)) {
            for (const json of row.items) {
                const item = ShopItem.fromJSON(json);
                if (item.item == null || item.item.itemId !== itemId || item.bundles === 0) {
                    continue;
                }
                entries.push({
                    ownerName: row.owner_name,
                    channelId: row.channel_id ?? -1,
                    mapId: row.map_id ?? 0,
                    title: row.title,
                    perBundle: item.perBundle,
                    bundles: item.bundles,
                    price: item.price,
                    sn: row.sn ?? 0,
                    item: item.item,
                });
            }
        }
        if (entries.length > 0) {
            await this.ranking.add(worldId, itemId);
        }
        entries.sort((a, b) => {
            if (highFirst) {
                return b.price - a.price;
            }
            return a.price - b.price;
        });
        return entries.slice(0, SearchLimit);
    }

    async findPopularShopSearches(worldId: number): Promise<number[]> {
        return this.ranking.top(worldId, PopularLimit);
    }
}
