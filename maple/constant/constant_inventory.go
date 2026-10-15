package constant

type InventoryType int8

const (
	InventoryTypeEquipment    InventoryType = 1
	InventoryTypeConsume      InventoryType = 2
	InventoryTypeInstallation InventoryType = 3
	InventoryTypeETC          InventoryType = 4
	InventoryTypeCash         InventoryType = 5
)
