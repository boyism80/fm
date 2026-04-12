package constant

type InventoryType int8

const (
	INVENTORY_TYPE_EQUIPMENT    InventoryType = 1
	INVENTORY_TYPE_CONSUME      InventoryType = 2
	INVENTORY_TYPE_INSTALLATION InventoryType = 3
	INVENTORY_TYPE_ETC          InventoryType = 4
	INVENTORY_TYPE_CASH         InventoryType = 5
)

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
