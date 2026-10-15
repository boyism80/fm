package game

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *InternalSystem) SearchShopsAsync(ctx actor.Context, ch *Character, itemID uint32, highFirst bool) *async.Promise[*internal.SearchShopsReply] {
	req := &internal.SearchShopsRequest{
		WorldId:   s.gw.GetWorldID(),
		ItemId:    itemID,
		HighFirst: highFirst,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SearchShopsReply, error) {
			return s.client.SearchShops(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) FindPopularShopSearchesAsync(ctx actor.Context, ch *Character) *async.Promise[*internal.FindPopularShopSearchesReply] {
	req := &internal.FindPopularShopSearchesRequest{WorldId: s.gw.GetWorldID()}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FindPopularShopSearchesReply, error) {
			return s.client.FindPopularShopSearches(c, req)
		},
		nil,
	)
}
