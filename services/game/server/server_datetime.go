package server

import (
	"context"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (gs *GameServer) SetServerDateTimeAsync(actorCtx actor.Context, reset bool, datetime string) *async.Promise {
	if gs == nil || gs.internalClient == nil {
		p := async.NewDeferred()
		p.SetError(fmt.Errorf("internal client unavailable"))
		return p
	}
	return async.ThenRPC(
		async.NewPromise(actorCtx, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.SetServerDateTimeReply, error) {
			return gs.internalClient.SetServerDateTime(c, &internal.SetServerDateTimeRequest{
				WorldId:  gs.config.WorldId,
				Datetime: datetime,
				Reset_:   reset,
			})
		},
		func(reply *internal.SetServerDateTimeReply) error {
			if reply == nil || !reply.GetOk() {
				return fmt.Errorf("set server datetime failed")
			}
			return nil
		},
	)
}
