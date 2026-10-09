package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (r *shopRoom) ToProto() *internal.Shop {
	pb := &internal.Shop{
		ShopId:         r.ID,
		AccountId:      r.AccountID,
		CharacterId:    r.OwnerID,
		OwnerName:      r.OwnerName,
		MapId:          r.MapID,
		ItemId:         r.ItemID,
		Title:          r.Title,
		Meso:           r.Meso,
		Items:          r.itemsToProto(r.Items),
		OpenedAtUnixMs: r.OpenedAt.UnixMilli(),
		Kind:           uint32(r.kind),
	}
	if r.published {
		pb.Sn = r.OID
	}
	for _, sale := range r.Sold {
		pb.Sold = append(pb.Sold, &internal.ShopSale{
			ItemId:  sale.ItemID,
			Bundles: uint32(sale.Bundles),
			Total:   sale.Total,
			Buyer:   sale.Buyer,
		})
	}
	return pb
}

func (r *shopRoom) storeBankToProto(meso int32, items []*ShopItem) *internal.Shop {
	return &internal.Shop{
		AccountId:   r.AccountID,
		CharacterId: r.OwnerID,
		OwnerName:   r.OwnerName,
		ItemId:      r.ItemID,
		Title:       r.Title,
		Meso:        meso,
		Items:       r.itemsToProto(items),
		Kind:        uint32(r.kind),
		StoreBankId: r.GameWorld.NewUniqueID(),
	}
}

func (r *shopRoom) itemsToProto(items []*ShopItem) []*internal.ShopItem {
	pb := make([]*internal.ShopItem, 0, len(items))
	for _, listed := range items {
		pb = append(pb, &internal.ShopItem{
			Item:      listed.Item.ToProto(0, 0),
			Bundles:   uint32(listed.Bundles),
			PerBundle: uint32(listed.PerBundle),
			Price:     listed.Price,
		})
	}
	return pb
}
