package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (hm *HiredMerchant) ToProto() *internal.HiredMerchant {
	pb := &internal.HiredMerchant{
		MerchantId:     hm.ID,
		AccountId:      hm.AccountID,
		CharacterId:    hm.OwnerID,
		OwnerName:      hm.OwnerName,
		MapId:          hm.MapID,
		ItemId:         hm.ItemID,
		Title:          hm.Title,
		Meso:           hm.Meso,
		OpenedAtUnixMs: hm.OpenedAt.UnixMilli(),
	}
	for _, listed := range hm.Items {
		pb.Items = append(pb.Items, &internal.HiredMerchantItem{
			Item:      listed.Item.ToProto(0, 0),
			Bundles:   uint32(listed.Bundles),
			PerBundle: uint32(listed.PerBundle),
			Price:     listed.Price,
		})
	}
	for _, sale := range hm.Sold {
		pb.Sold = append(pb.Sold, &internal.HiredMerchantSale{
			ItemId:  sale.ItemID,
			Bundles: uint32(sale.Bundles),
			Total:   sale.Total,
			Buyer:   sale.Buyer,
		})
	}
	return pb
}
