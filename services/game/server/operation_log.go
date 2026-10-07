package server

import (
	"context"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

func (gs *GameServer) WriteOperationLogAsync(ctx actor.Context, charID uint32, kind constant.OperationLogKind, meso int32, detail string) *async.Promise[*internal.WriteOperationLogReply] {
	log.Printf("operation log: kind=%s character=%d meso=%d detail=%s", kind, charID, meso, detail)
	if gs.internalClient == nil {
		return nil
	}
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	return promise.ThenRPC(func(c context.Context) (*internal.WriteOperationLogReply, error) {
		return gs.internalClient.WriteOperationLog(c, &internal.WriteOperationLogRequest{
			WorldId:     gs.config.WorldId,
			ChannelId:   uint32(gs.config.ChannelId),
			CharacterId: charID,
			Kind:        string(kind),
			Meso:        int64(meso),
			Detail:      detail,
		})
	}, func(reply *internal.WriteOperationLogReply) error {
		if reply.GetOk() == false {
			log.Printf("operation log: rejected kind=%s character=%d", kind, charID)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("operation log: kind=%s character=%d: %v", kind, charID, err)
	})
}
