package actorx

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
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
	handler   *handler.MessageHandler
	mapPID    *actor.PID
	scheduler *scheduler.TimerScheduler
}

func NewItemActor(ctx actor.Context,
	serverCtx *context.ServerContext,
	entity entity.Item,
	mapPID *actor.PID) actor.Actor {

	actor := &ItemActor{
		handler:   handler.NewMessageHandler(),
		Item:      entity,
		mapPID:    mapPID,
		scheduler: scheduler.NewTimerScheduler(ctx),
	}
	RegisterObjectHandlers(ctx, actor.Item.GetObject(), actor.handler)
	handler.RegisterHandler(ctx, actor, actor.handler, onItemStarted)
	handler.RegisterHandler(ctx, actor, actor.handler, onItemSpawn)
	handler.RegisterHandler(ctx, actor, actor.handler, onItemLooting)
	handler.RegisterHandler(ctx, actor, actor.handler, onItemLooted)
	handler.RegisterHandler(ctx, actor, actor.handler, onItemDropTypeChanged)
	handler.RegisterHandler(ctx, actor, actor.handler, onItemDestroy)

	return actor
}

func (state *ItemActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onItemStarted(ctx actor.Context, state *ItemActor, m *actor.Started) {
	state.scheduler.SendOnce(30*time.Second, ctx.Self(), &msg.ItemDropTypeChanged{
		Mode: constant.DROP_TYPE_FFA,
	})

	state.scheduler.SendOnce(2*time.Minute, ctx.Self(), &msg.ItemDestroy{
		Animation: constant.DROP_ITEM_ANIMATION_TYPE_DISAPPEAR,
	})
}

func onItemSpawn(ctx actor.Context, state *ItemActor, m *msg.Spawn) {
	drop := state.GetDrop()
	ctx.Send(state.mapPID, &msg.MapBroadcast{
		Sender:     ctx.Self(),
		Pivot:      state.GetObject().Position,
		ExceptSelf: true,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.SpawnItem{
				ID:           drop.ID,
				Animation:    constant.DROP_ITEM_ANIMATION_TYPE_LOOTING,
				DropType:     drop.DropType,
				Item:         state.Item,
				OwnerID:      drop.Owner,
				SpawnedPoint: drop.SpawnedPoint,
				IsPlayerDrop: true,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func onItemLooting(ctx actor.Context, state *ItemActor, m *msg.ItemLooting) {
	drop := state.GetDrop()
	if drop.Looting {
		ctx.Send(m.Actor, &msg.CharacterLootFailed{
			OID: drop.ID,
		})
		return
	}

	if drop.DropType == constant.DROP_TYPE_OWNED && drop.Owner != m.CharacterId {
		ctx.Send(m.Actor, &msg.CharacterLootFailed{
			OID: drop.ID,
		})
		return
	}

	ctx.Send(m.Actor, &msg.CharacterItemLooting{
		PID:  ctx.Self(),
		OID:  drop.ID,
		Item: state.Item.Clone(state.GetCount()),
	})
	drop.Looting = true
}

func onItemLooted(ctx actor.Context, state *ItemActor, m *msg.ItemLooted) {
	drop := state.GetDrop()
	drop.Looting = false
	if m.Success {
		if state.Reduce(uint16(m.Count)) == 0 {
			ctx.Send(state.mapPID, &msg.MapRemoveItem{
				Actor:       ctx.Self(),
				OID:         drop.ID,
				CharacterId: m.CharacterId,
				Mode:        resp.REMOVE_ITEM_TYPE_ANIMATED,
				Position:    drop.Position,
			})
		}
	}
}

func onItemDropTypeChanged(ctx actor.Context, state *ItemActor, m *msg.ItemDropTypeChanged) {
	drop := state.GetDrop()
	drop.DropType = m.Mode
}

func onItemDestroy(ctx actor.Context, state *ItemActor, m *msg.ItemDestroy) {
	drop := state.GetDrop()
	ctx.Send(state.mapPID, &msg.MapRemoveItem{
		Actor:    ctx.Self(),
		OID:      drop.ID,
		Mode:     resp.REMOVE_ITEM_TYPE_EXPIRED,
		Position: drop.Position,
	})
}
