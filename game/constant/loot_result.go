package constant

// LootResult defines the different results for loot failure or success
type LootResult uint8

const (
	LOOT_SUCCESS               LootResult = 0
	LOOT_FAILED_ITEM_NOT_FOUND LootResult = 1
	LOOT_FAILED_NO_OWNERSHIP   LootResult = 2
	LOOT_FAILED_INVENTORY_FULL LootResult = 3
	LOOT_FAILED_MESO_FULL      LootResult = 4
	LOOT_FAILED_INVALID_ITEM   LootResult = 5
)
