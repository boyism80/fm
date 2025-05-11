package constant

type ItemType int

const (
	ItemTypeEquipment ItemType = 1
	ItemTypeEtc       ItemType = 2
)

type EquipmentPartsType int

const (
	EquipmentPartsWeapon EquipmentPartsType = -11
	EquipmentPartsShield EquipmentPartsType = -10
)

type InventoryType uint8

const (
	InventoryTypeEquip        InventoryType = 1
	InventoryTypeConsume      InventoryType = 2
	InventoryTypeInstallation InventoryType = 3
	InventoryTypeEtc          InventoryType = 4
	InventoryTypeCash         InventoryType = 5
)
