package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (es *EntrustedShop) ToProto() *internal.EntrustedShop {
	pb := &internal.EntrustedShop{
		ShopId:         es.ID,
		AccountId:      es.AccountID,
		CharacterId:    es.OwnerID,
		OwnerName:      es.OwnerName,
		MapId:          es.MapID,
		ItemId:         es.ItemID,
		Title:          es.Title,
		Meso:           es.Meso,
		OpenedAtUnixMs: es.OpenedAt.UnixMilli(),
	}
	for _, listed := range es.Items {
		pb.Items = append(pb.Items, &internal.EntrustedShopItem{
			Item:      listed.Item.ToProto(0, 0),
			Bundles:   uint32(listed.Bundles),
			PerBundle: uint32(listed.PerBundle),
			Price:     listed.Price,
		})
	}
	for _, sale := range es.Sold {
		pb.Sold = append(pb.Sold, &internal.EntrustedShopSale{
			ItemId:  sale.ItemID,
			Bundles: uint32(sale.Bundles),
			Total:   sale.Total,
			Buyer:   sale.Buyer,
		})
	}
	return pb
}
