package entity

type Item struct {
	Object
	Id            uint32
	UniqueID      uint64
	Type          uint8
	Parts         int16
	Quantity      uint16
	Owner         string
	Flag          uint16
	Expiration    int64
	Pet           *Pet
	InventoryID   uint64
	UpgradeSlots  uint8
	Level         uint8
	Str           uint16
	Dex           uint16
	Int           uint16
	Luk           uint16
	Hp            uint16
	Mp            uint16
	Watk          uint16
	Matk          uint16
	Wdef          uint16
	Mdef          uint16
	Acc           uint16
	Avoid         uint16
	Hands         uint16
	Speed         uint16
	Jump          uint16
	IncSkill      uint16
	BaseLevel     uint8
	EquipLevel    uint8
	ExpPercentage uint32
	Equip         *Equip
}

func (i *Item) IsThrowingStar() bool {
	return i.Id/10000 == 207
}

func (i *Item) IsBullet() bool {
	return i.Id/10000 == 233
}
