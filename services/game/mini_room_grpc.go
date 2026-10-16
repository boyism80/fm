package game

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (r *shopRoom) BuildRPC() *internal.Shop {
	pb := &internal.Shop{
		ShopId:         r.id,
		AccountId:      r.accountID,
		CharacterId:    r.ownerID,
		OwnerName:      r.ownerName,
		MapId:          r.mapID,
		ItemId:         r.itemID,
		Title:          r.title,
		Meso:           r.meso,
		Items:          r.buildItemsRPC(r.items),
		OpenedAtUnixMs: r.openedAt.UnixMilli(),
		Kind:           uint32(r.kind),
	}
	if r.published {
		pb.Sn = r.oid
	}
	for _, sale := range r.sold {
		pb.Sold = append(pb.Sold, &internal.ShopSale{
			ItemId:  sale.ItemID,
			Bundles: uint32(sale.Bundles),
			Total:   sale.Total,
			Buyer:   sale.Buyer,
		})
	}
	return pb
}

func (r *shopRoom) buildStoreBank(meso int32, items []*ShopItem) closedShop {
	return closedShop{
		id:    r.gameWorld.NewUniqueID(),
		meso:  meso,
		items: items,
	}
}

func (s closedShop) BuildRPC(r *shopRoom) *internal.Shop {
	return &internal.Shop{
		AccountId:   r.accountID,
		CharacterId: r.ownerID,
		OwnerName:   r.ownerName,
		ItemId:      r.itemID,
		Title:       r.title,
		Meso:        s.meso,
		Items:       r.buildItemsRPC(s.items),
		Kind:        uint32(r.kind),
		StoreBankId: s.id,
	}
}

func (r *shopRoom) buildItemsRPC(items []*ShopItem) []*internal.ShopItem {
	pb := make([]*internal.ShopItem, 0, len(items))
	for _, listed := range items {
		pb = append(pb, &internal.ShopItem{
			Item:      listed.item.BuildRPC(0, 0),
			Bundles:   uint32(listed.bundles),
			PerBundle: uint32(listed.perBundle),
			Price:     listed.price,
		})
	}
	return pb
}
