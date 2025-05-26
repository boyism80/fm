package actor

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
	mapPid    *actor.PID
	scheduler *scheduler.TimerScheduler
}

func NewItemActor(ctx actor.Context,
	serverCtx *context.ServerContext,
	entity entity.Item,
	mapPid *actor.PID) actor.Actor {

	actor := &ItemActor{
		handler:   handler.NewMessageHandler(),
		Item:      entity,
		mapPid:    mapPid,
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
		Mode: constant.DropTypeFFA,
	})

	state.scheduler.SendOnce(2*time.Minute, ctx.Self(), &msg.ItemDestroy{
		Animation: constant.DropItemAnimationTypeDisappear,
	})
}

func onItemSpawn(ctx actor.Context, state *ItemActor, m *msg.ItemSpawn) {
	drop := state.GetDrop()
	ctx.Send(state.mapPid, &msg.MapBroadcastRange{
		Sender:     ctx.Self(),
		Pivot:      state.GetObject().Position,
		ExceptSelf: true,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.DropItem{
				Id:           drop.Id,
				Animation:    constant.DropItemAnimationTypeLooting,
				DropType:     drop.DropType,
				Item:         state.Item,
				OwnerId:      drop.Owner,
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
			Oid: drop.Id,
		})
		return
	}

	if drop.DropType == constant.DropTypeOwned && drop.Owner != m.CharacterId {
		ctx.Send(m.Actor, &msg.CharacterLootFailed{
			Oid: drop.Id,
		})
		return
	}

	ctx.Send(m.Actor, &msg.CharacterItemLooting{
		Pid:  ctx.Self(),
		Oid:  drop.Id,
		Item: state.Item.Clone(state.GetCount()),
	})
	drop.Looting = true
}

func onItemLooted(ctx actor.Context, state *ItemActor, m *msg.ItemLooted) {
	drop := state.GetDrop()
	drop.Looting = false
	if m.Success {
		if state.Reduce(uint16(m.Count)) == 0 {
			ctx.Send(state.mapPid, &msg.MapRemoveItem{
				Actor:       ctx.Self(),
				Oid:         drop.Id,
				CharacterId: m.CharacterId,
				Mode:        resp.RemoveItemTypeAnimated,
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
	ctx.Send(state.mapPid, &msg.MapRemoveItem{
		Actor:    ctx.Self(),
		Oid:      drop.Id,
		Mode:     resp.RemoveItemTypeExpired,
		Position: drop.Position,
	})
}
