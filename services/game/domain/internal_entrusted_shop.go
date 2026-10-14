package domain

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *InternalSystem) FindEntrustedShopAsync(ctx actor.Context, ch *Character) *async.Promise[*internal.FindEntrustedShopReply] {
	req := &internal.FindEntrustedShopRequest{
		WorldId:   s.gw.GetWorldID(),
		AccountId: ch.AccountID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.FindEntrustedShopReply, error) {
			return s.client.FindEntrustedShop(c, req)
		},
		nil,
	)
}

func (s *InternalSystem) OpenShopAsync(ctx actor.Context, ch *Character, shop *internal.Shop) *async.Promise[*internal.OpenShopReply] {
	shop.WorldId = s.gw.GetWorldID()
	shop.ChannelId = s.gw.GetChannelID()
	req := &internal.OpenShopRequest{Shop: shop}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.OpenShopReply, error) {
			return s.client.OpenShop(c, req)
		},
		nil,
	)
}
