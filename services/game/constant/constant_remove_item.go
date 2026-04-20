package constant

type RemoveItemType uint8

const (
	REMOVE_ITEM_TYPE_EXPIRED     RemoveItemType = 0
	REMOVE_ITEM_TYPE_NO_ANIMATED RemoveItemType = 1
	REMOVE_ITEM_TYPE_ANIMATED    RemoveItemType = 2
	REMOVE_ITEM_TYPE_EXPLOSION   RemoveItemType = 4
	REMOVE_ITEM_TYPE_LOOT_BY_PET RemoveItemType = 5
)
