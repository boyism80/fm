package actor

import (
	"log"

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
	looting bool
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
	handler.RegisterHandler(ctx, act, act.handler, onItemLooting)
	handler.RegisterHandler(ctx, act, act.handler, onItemLooted)
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

func onItemLooting(ctx actor.Context, state *ItemActor, request *msg.ItemLooting) {
	if state.looting {
		ctx.Send(request.Actor, &msg.CharacterLootFailed{
			Oid: state.id,
		})
		return
	}

	state.looting = true
	ctx.Send(request.Actor, &msg.CharacterItemLooting{
		Pid:  ctx.Self(),
		Oid:  state.id,
		Item: state.Item.Clone(state.GetCount()),
	})
}

func onItemLooted(ctx actor.Context, state *ItemActor, request *msg.ItemLooted) {
	state.looting = false
	if request.Success {
		if state.Reduce(request.Count) == 0 {
			// delete item
			// send map for remove me

			ctx.Send(state.mapPid, &msg.MapItemLooted{
				Actor:       ctx.Self(),
				Oid:         state.id,
				CharacterId: request.CharacterId,
				Mode:        resp.RemoveItemTypeAnimated,
				Position:    state.GetObject().Position,
			})
		}
	}

	log.Println("item looted")
}
