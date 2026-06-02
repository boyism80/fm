package constant

type LootResult uint8

const (
	LootSuccess             LootResult = 0
	LootFailedItemNotFound  LootResult = 1
	LootFailedNoOwnership   LootResult = 2
	LootFailedInventoryFull LootResult = 3
	LootFailedMesoFull      LootResult = 4
	LootFailedInvalidItem   LootResult = 5
)
