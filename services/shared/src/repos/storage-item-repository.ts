import { InventoryRepository } from "./inventory-repository";

export class StorageItemRepository extends InventoryRepository {
    protected override get table() {
        return "storage_items";
    }
}
