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
		OpenedAtUnixMs: r.OpenedAt.UnixMilli(),
		Kind:           uint32(r.kind),
	}
	if r.published {
		pb.Sn = r.OID
	}
	for _, listed := range r.Items {
		pb.Items = append(pb.Items, &internal.ShopItem{
			Item:      listed.Item.ToProto(0, 0),
			Bundles:   uint32(listed.Bundles),
			PerBundle: uint32(listed.PerBundle),
			Price:     listed.Price,
		})
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
