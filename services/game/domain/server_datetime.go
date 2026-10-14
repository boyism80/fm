package domain

import (
	"context"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *InternalSystem) SetServerDateTimeAsync(actorCtx actor.Context, reset bool, datetime string) *async.Promise[*internal.SetServerDateTimeReply] {
	return async.NewTask(actorCtx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SetServerDateTimeReply, error) {
			return s.client.SetServerDateTime(c, &internal.SetServerDateTimeRequest{
				WorldId:  s.gw.GetWorldID(),
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
