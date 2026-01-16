package wz

// ShopItem represents a single item in an NPC shop
type ShopItem struct {
	ItemID   uint32  // Item ID
	Price    int     // Price in mesos (0 if unitPrice is used)
	Period   int     // Period/reset cycle
	Stock    int     // Stock quantity
	UnitPrice float64 // Unit price for throwing stars/bullets (0 if price is used)
}

// Shop represents an NPC shop with its items
type Shop struct {
	NpcID uint32      // NPC ID (shop owner)
	Items []ShopItem  // List of items in the shop
}
