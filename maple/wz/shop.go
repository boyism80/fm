package wz

type ShopItem struct {
	ItemID uint32
	Price  int
}

type Shop struct {
	NpcID uint32
	Items []ShopItem
}

func (s *Shop) Find(itemID uint32) *ShopItem {
	for i := range s.Items {
		if s.Items[i].ItemID == itemID {
			return &s.Items[i]
		}
	}
	return nil
}
