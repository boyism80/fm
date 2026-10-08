package server

import (
	"context"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) SearchShopsAsync(ctx actor.Context, ch *entity.Character, itemID uint32, highFirst bool) *async.Promise[*internal.SearchShopsReply] {
	req := &internal.SearchShopsRequest{
		WorldId:   l.gs.config.WorldId,
		ItemId:    itemID,
		HighFirst: highFirst,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SearchShopsReply, error) {
			return l.gs.internalClient.SearchShops(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) FindPopularShopSearchesAsync(ctx actor.Context, ch *entity.Character) *async.Promise[*internal.FindPopularShopSearchesReply] {
	req := &internal.FindPopularShopSearchesRequest{WorldId: l.gs.config.WorldId}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FindPopularShopSearchesReply, error) {
			return l.gs.internalClient.FindPopularShopSearches(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) OnShopScannerResult(ch *entity.Character, itemID uint32, entries []*internal.ShopSearchEntry) {
	result := &response.ShopScannerResult{ItemID: itemID}
	for _, entry := range entries {
		item, err := entity.NewItemFromInternalProto(entry.GetItem(), ch.GameWorld)
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
			listed.Item = item.ToDTO()
		}
		result.Entries = append(result.Entries, listed)
	}
	ch.Send(result, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShopScannerPopular(ch *entity.Character, itemIDs []uint32) {
	ch.Send(&response.ShopScannerPopular{ItemIDs: itemIDs}, types.SEND_POLICY_ENCRYPT)
}
