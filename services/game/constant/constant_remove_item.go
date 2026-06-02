package constant

type RemoveItemType uint8

const (
	RemoveItemTypeExpired    RemoveItemType = 0
	RemoveItemTypeNoAnimated RemoveItemType = 1
	RemoveItemTypeAnimated   RemoveItemType = 2
	RemoveItemTypeExplosion  RemoveItemType = 4
	RemoveItemTypeLootByPet  RemoveItemType = 5
)
