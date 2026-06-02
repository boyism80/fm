package constant

type InventoryType int8

const (
	InventoryTypeEquipment    InventoryType = 1
	InventoryTypeConsume      InventoryType = 2
	InventoryTypeInstallation InventoryType = 3
	InventoryTypeETC          InventoryType = 4
	InventoryTypeCash         InventoryType = 5
)

func GetInventoryTypeByItemID(itemID uint32) InventoryType {
	itemType := itemID / 10000
	switch {
	case itemType >= 100 && itemType < 200:
		return InventoryTypeEquipment
	case itemType >= 200 && itemType < 300:
		return InventoryTypeConsume
	case itemType >= 300 && itemType < 400:
		return InventoryTypeInstallation
	case itemType >= 400 && itemType < 500:
		return InventoryTypeETC
	case itemType >= 500 && itemType < 600:
		return InventoryTypeCash
	default:
		return InventoryTypeETC
	}
}
