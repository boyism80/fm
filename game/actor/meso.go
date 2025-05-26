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

type MesoActor struct {
	entity.Meso
	handler   *handler.MessageHandler
	mapPid    *actor.PID
	scheduler *scheduler.TimerScheduler
}

func NewMesoActor(ctx actor.Context,
	serverCtx *context.ServerContext,
	meso entity.Meso,
	mapPid *actor.PID,
	position types.Vector2[int16]) actor.Actor {

	actor := &MesoActor{
		Meso:      meso,
		handler:   handler.NewMessageHandler(),
		mapPid:    mapPid,
		scheduler: scheduler.NewTimerScheduler(ctx),
	}
	RegisterObjectHandlers(ctx, actor.Drop.Object, actor.handler)
	handler.RegisterHandler(ctx, actor, actor.handler, onMesoStarted)
	handler.RegisterHandler(ctx, actor, actor.handler, onMesoSpawn)
	handler.RegisterHandler(ctx, actor, actor.handler, onMesoLooting)
	handler.RegisterHandler(ctx, actor, actor.handler, onMesoLooted)
	handler.RegisterHandler(ctx, actor, actor.handler, onMesoDropTypeChanged)
	handler.RegisterHandler(ctx, actor, actor.handler, onMesoDestroy)
	return actor
}

func (state *MesoActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onMesoStarted(ctx actor.Context, state *MesoActor, m *actor.Started) {
	state.scheduler.SendOnce(30*time.Second, ctx.Self(), &msg.ItemDropTypeChanged{
		Mode: constant.DropTypeFFA,
	})

	state.scheduler.SendOnce(2*time.Minute, ctx.Self(), &msg.ItemDestroy{
		Animation: constant.DropItemAnimationTypeDisappear,
	})
}

func onMesoSpawn(ctx actor.Context, state *MesoActor, m *msg.MesoSpawn) {
	drop := state.Drop

	ctx.Send(state.mapPid, &msg.MapBroadcastRange{
		Sender:     ctx.Self(),
		Pivot:      state.Object.Position,
		ExceptSelf: true,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.DropMeso{
				Id:           drop.Id,
				Animation:    constant.DropItemAnimationTypeLooting,
				DropType:     drop.DropType,
				Count:        state.Count,
				OwnerId:      drop.Owner,
				Position:     drop.Position,
				SpawnedPoint: drop.SpawnedPoint,
				IsPlayerDrop: true,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func onMesoLooting(ctx actor.Context, state *MesoActor, m *msg.ItemLooting) {
	drop := state.Drop
	if state.Looting {
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

	ctx.Send(m.Actor, &msg.CharacterMesoLooting{
		Pid:  ctx.Self(),
		Oid:  drop.Id,
		Meso: state.Meso.Count,
	})
	drop.Looting = true
}

func onMesoLooted(ctx actor.Context, state *MesoActor, m *msg.ItemLooted) {
	drop := state.Drop
	state.Looting = false
	if m.Success {
		state.Count -= m.Count
		if state.Count == 0 {
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

func onMesoDropTypeChanged(ctx actor.Context, state *MesoActor, m *msg.ItemDropTypeChanged) {
	state.Drop.DropType = m.Mode
}

func onMesoDestroy(ctx actor.Context, state *MesoActor, m *msg.ItemDestroy) {
	drop := state.Drop
	ctx.Send(state.mapPid, &msg.MapRemoveItem{
		Actor:    ctx.Self(),
		Oid:      drop.Id,
		Mode:     resp.RemoveItemTypeExpired,
		Position: drop.Position,
	})
}
