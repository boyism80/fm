package server

import (
	"log"

	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnShopScannerResult(ch *game.Character, itemID uint32, entries []*internal.ShopSearchEntry) {
	result := &response.ShopScannerResult{ItemID: itemID}
	for _, entry := range entries {
		item, err := game.NewItemFromInternalProto(entry.GetItem(), ch.GameWorld)
		if err != nil {
			log.Printf("OnShopScannerResult item=%d: %v", itemID, err)
			continue
		}
		listed := response.ShopScannerEntry{
			OwnerName:     entry.GetOwnerName(),
			MapID:         entry.GetMapId(),
			Title:         entry.GetTitle(),
			PerBundle:     entry.GetPerBundle(),
			Bundles:       entry.GetBundles(),
			Price:         entry.GetPrice(),
			SN:            entry.GetSn(),
			Channel:       int8(entry.GetChannelId()),
			InventoryType: uint8(item.GetInventoryType()),
		}
		if item.GetInventoryType() == constant.InventoryTypeEquipment {
			listed.Item = item.BuildS2C()
		}
		result.Entries = append(result.Entries, listed)
	}
	ch.Send(result, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShopScannerPopular(ch *game.Character, itemIDs []uint32) {
	ch.Send(&response.ShopScannerPopular{ItemIDs: itemIDs}, types.SEND_POLICY_ENCRYPT)
}
