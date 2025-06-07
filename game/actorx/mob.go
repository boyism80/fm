package actorx

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

type MobActor struct {
	entity.Mob
	handler *handler.MessageHandler
	mapPID  *actor.PID
}

func NewMobActor(ctx actor.Context, serverCtx context.ServerContext, entity entity.Mob, mapPID *actor.PID) actor.Actor {
	actor := &MobActor{
		Mob:     entity,
		handler: handler.NewMessageHandler(),
		mapPID:  mapPID,
	}
	RegisterLifeHandlers(ctx, &actor.Mob.Life, actor.handler)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobSpawn)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobPlayerWarped)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobKill)
	return actor
}

func (state *MobActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onMobSpawn(ctx actor.Context, state *MobActor, m *msg.Spawn) {
	ctx.Send(state.mapPID, &msg.MapBroadcastRange{
		Sender:     ctx.Self(),
		Pivot:      state.Position,
		ExceptSelf: true,
		Message: &common_msg.SendProtocol{
			Protocol: &resp.SpawnMob{
				Mob:       &state.Mob,
				SpawnType: constant.MobSpawnTypeAnimate,
				Link:      0,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}

func onMobPlayerWarped(ctx actor.Context, state *MobActor, m *msg.Warped) {
	ctx.Send(m.Sender, &common_msg.SendProtocol{
		Protocol: &resp.SpawnMob{
			Mob:       &state.Mob,
			SpawnType: constant.MobSpawnTypeNone,
			Link:      0,
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})
}

func onMobKill(ctx actor.Context, state *MobActor, m *msg.MobKill) {
	ctx.Send(state.mapPID, &msg.MapDieMob{
		Sender:        ctx.Self(),
		OID:           state.ID,
		AnimationType: m.AnimationType,
		Position:      state.Position,
	})
}
