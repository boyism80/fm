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
	DropItemAnimationTypeLooting       DropItemAnimationType = 1
	DropItemAnimationTypeNone          DropItemAnimationType = 2
	DropItemAnimationTypeDisappearFade DropItemAnimationType = 3
	DropItemAnimationTypeDisappear     DropItemAnimationType = 4
)

type DropType uint8

const (
	DropTypeFFA   DropType = 0
	DropTypeParty DropType = 1
	DropTypeOwned DropType = 2
)

type Stat uint32

const (
	StatSkin        Stat = 0x1
	StatFace        Stat = 0x2
	StatHair        Stat = 0x4
	StatPet         Stat = 0x8
	StatLevel       Stat = 0x10
	StatJob         Stat = 0x20
	StatStr         Stat = 0x40
	StatDex         Stat = 0x80
	StatInt         Stat = 0x100
	StatLuk         Stat = 0x200
	StatHP          Stat = 0x400
	StatMaxHP       Stat = 0x800
	StatMP          Stat = 0x1000
	StatMaxMP       Stat = 0x2000
	StatAvailableAP Stat = 0x4000
	StatAvailableSP Stat = 0x8000
	StatExp         Stat = 0x10000
	StatFame        Stat = 0x20000
	StatMeso        Stat = 0x40000
)
