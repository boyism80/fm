package constant

// InventoryType represents the different inventory tabs.
type InventoryType int8

const (
	INVENTORY_TYPE_EQUIPMENT    InventoryType = 1 // Equipment inventory
	INVENTORY_TYPE_CONSUME      InventoryType = 2 // Consumable items
	INVENTORY_TYPE_INSTALLATION InventoryType = 3 // Installation items
	INVENTORY_TYPE_ETC          InventoryType = 4 // Miscellaneous items
	INVENTORY_TYPE_CASH         InventoryType = 5 // Cash shop items
)

// GetInventoryTypeByItemID returns the inventory type for an item by its ID.
// Uses MapleStory item ID ranges (itemID / 10000).
func GetInventoryTypeByItemID(itemID uint32) InventoryType {
	itemType := itemID / 10000
	switch {
	case itemType >= 100 && itemType < 200:
		return INVENTORY_TYPE_EQUIPMENT
	case itemType >= 200 && itemType < 300:
		return INVENTORY_TYPE_CONSUME
	case itemType >= 300 && itemType < 400:
		return INVENTORY_TYPE_INSTALLATION
	case itemType >= 400 && itemType < 500:
		return INVENTORY_TYPE_ETC
	case itemType >= 500 && itemType < 600:
		return INVENTORY_TYPE_CASH
	default:
		return INVENTORY_TYPE_ETC
	}
}
