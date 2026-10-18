package game

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *InternalSystem) GetRankingAsync(ctx actor.Context, kind internal.RankingKind, count uint32) *async.Promise[*internal.GetRankingReply] {
	req := &internal.GetRankingRequest{
		Kind:    kind,
		WorldId: s.gw.GetWorldID(),
		Count:   count,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.GetRankingReply, error) {
			return s.client.GetRanking(c, req)
		},
		nil,
	)
}
