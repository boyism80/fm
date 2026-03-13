package constant

// RemoveItemType represents how an item is removed from the map.
type RemoveItemType uint8

const (
	REMOVE_ITEM_TYPE_EXPIRED RemoveItemType = iota
	REMOVE_ITEM_TYPE_NO_ANIMATED
	REMOVE_ITEM_TYPE_ANIMATED
	REMOVE_ITEM_TYPE_EXPLOSION
	REMOVE_ITEM_TYPE_LOOT_BY_PET
)
