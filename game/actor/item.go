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
	oid     uint32
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
		oid:     oid,
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
				Id:           state.oid,
				Animation:    constant.DropItemAnimationTypeDefault,
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
			Oid: state.oid,
		})
	} else {
		ctx.Send(request.Actor, &msg.CharacterItemLooting{
			Pid:  ctx.Self(),
			Oid:  state.oid,
			Item: state.Item.Clone(state.GetCount()),
		})
		state.looting = true
	}
}

func onItemLooted(ctx actor.Context, state *ItemActor, request *msg.ItemLooted) {
	state.looting = false
	if request.Success {
		if state.Reduce(uint16(request.Count)) == 0 {
			ctx.Send(state.mapPid, &msg.MapItemLooted{
				Actor:       ctx.Self(),
				Oid:         state.oid,
				CharacterId: request.CharacterId,
				Mode:        resp.RemoveItemTypeAnimated,
				Position:    state.GetObject().Position,
			})
		}
	}

	log.Println("item looted")
}
