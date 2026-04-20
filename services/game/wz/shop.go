package wz

type ShopItem struct {
	ItemID    uint32
	Price     int
	Period    int
	Stock     int
	UnitPrice float64
}

type Shop struct {
	NpcID uint32
	Items []ShopItem
}
