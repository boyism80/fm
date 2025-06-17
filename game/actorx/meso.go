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

type MesoActor struct {
	entity.Meso
	handler   *handler.MessageHandler
	mapPID    *actor.PID
	scheduler *scheduler.TimerScheduler
}

func NewMesoActor(ctx actor.Context,
	serverCtx *context.ServerContext,
	meso entity.Meso,
	mapPID *actor.PID,
	position types.Vector2[int16]) actor.Actor {

	actor := &MesoActor{
		Meso:      meso,
		handler:   handler.NewMessageHandler(),
		mapPID:    mapPID,
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
		Mode: constant.DROP_TYPE_FFA,
	})

	state.scheduler.SendOnce(2*time.Minute, ctx.Self(), &msg.ItemDestroy{
		Animation: constant.DROP_ITEM_ANIMATION_TYPE_DISAPPEAR,
	})
}

func onMesoSpawn(ctx actor.Context, state *MesoActor, m *msg.Spawn) {
	drop := state.Drop

	ctx.Send(state.mapPID, &msg.MapBroadcast{
		Sender:     ctx.Self(),
		Pivot:      state.Object.Position,
		ExceptSelf: true,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.SpawnMeso{
				ID:           drop.ID,
				Animation:    constant.DROP_ITEM_ANIMATION_TYPE_LOOTING,
				DropType:     drop.DropType,
				Count:        state.Count,
				OwnerID:      drop.Owner,
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

	ctx.Send(m.Actor, &msg.CharacterMesoLooting{
		PID:  ctx.Self(),
		OID:  drop.ID,
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

func onMesoDropTypeChanged(ctx actor.Context, state *MesoActor, m *msg.ItemDropTypeChanged) {
	state.Drop.DropType = m.Mode
}

func onMesoDestroy(ctx actor.Context, state *MesoActor, m *msg.ItemDestroy) {
	drop := state.Drop
	ctx.Send(state.mapPID, &msg.MapRemoveItem{
		Actor:    ctx.Self(),
		OID:      drop.ID,
		Mode:     resp.REMOVE_ITEM_TYPE_EXPIRED,
		Position: drop.Position,
	})
}
