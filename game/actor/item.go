package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/handler"
	common_msg "github.com/boyism80/fm/common/msg"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
	"github.com/boyism80/fm/game/protocol/resp"
)

type ItemActor struct {
	entity.Item
	id      uint32
	ctx     *context.ServerContext
	handler *handler.MessageHandler
	mapPid  *actor.PID
}

func NewItemActor(ctx actor.Context,
	serverCtx *context.ServerContext,
	entity entity.Item,
	oid uint32,
	mapPid *actor.PID) actor.Actor {

	act := &ItemActor{
		id:      oid,
		ctx:     serverCtx,
		handler: handler.NewMessageHandler(),
		Item:    entity,
		mapPid:  mapPid,
	}
	RegisterObjectHandlers(ctx, act.Item.GetObject(), act.handler)
	handler.RegisterHandler(ctx, act, act.handler, onItemSpawn)
	return act
}

func (state *ItemActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onItemSpawn(ctx actor.Context, state *ItemActor, request *msg.ItemSpawn) {
	ctx.Send(state.mapPid, &msg.MapBroadcastRange{
		Sender:     ctx.Self(),
		Pivot:      state.GetObject().Position,
		ExceptSelf: true,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.DropItem{
				Id:           state.id,
				Animation:    constant.DropItemAnimationTypeDefault,
				Meso:         0,
				DropType:     2,
				Item:         state.Item,
				OwnerId:      request.OwnerId,
				SpawnedPoint: request.SpawnedPoint,
				IsPlayerDrop: true,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}
