package constant

type ItemType int

const (
	ItemTypeEquipment ItemType = 1
	ItemTypeEtc       ItemType = 2
)

type EquipmentPartsType int16

const (
	EquipmentPartsWeapon EquipmentPartsType = -11
	EquipmentPartsShield EquipmentPartsType = -10
)

type InventoryType int8

const (
	InventoryTypeEquipment    InventoryType = 1
	InventoryTypeConsume      InventoryType = 2
	InventoryTypeInstallation InventoryType = 3
	InventoryTypeEtc          InventoryType = 4
	InventoryTypeCash         InventoryType = 5
)

type DropItemAnimationType uint8

const (
	DropItemAnimationTypeDefault       DropItemAnimationType = 1
	DropItemAnimationTypeNone          DropItemAnimationType = 2
	DropItemAnimationTypeDisappearFade DropItemAnimationType = 3
	DropItemAnimationTypeDisappear     DropItemAnimationType = 4
)
