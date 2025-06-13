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
	"github.com/boyism80/fm/game/protocol"
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
	handler.RegisterHandler(ctx, actor, actor.handler, onMobControllerChange)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobMove)
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
		OID:           state.OID,
		AnimationType: m.AnimationType,
		Position:      state.Position,
	})
}

func onMobControllerChange(ctx actor.Context, state *MobActor, m *msg.MobControllerChange) {
	ctx.Send(m.Controller, &common_msg.SendProtocol{
		Protocol: &resp.ControlMob{
			Mob:   &state.Mob,
			Aggro: false,
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})
}

func onMobMove(ctx actor.Context, state *MobActor, m *msg.MobMove) {
	ctx.Send(m.Sender, &common_msg.SendProtocol{
		Protocol: &resp.ControlMoveMob{
			OID:          state.OID,
			MoveId:       m.MovementId,
			EnabledSkill: m.IsAggroed,
			MP:           state.Mp,
			SkillId:      0,
			SkillLevel:   0,
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})

	for _, movement := range m.Movements {
		if move, ok := movement.(*protocol.AbsoluteLifeMovement); ok {
			state.Mob.Position = move.Position
		}

		state.Mob.Stance = movement.GetStance()
	}

	ctx.Send(state.mapPID, &msg.MapBroadcastRange{
		Sender:     ctx.Self(),
		Pivot:      state.Position,
		ExceptSelf: true,
		Excepts:    map[*actor.PID]struct{}{m.Sender: {}},
		Message: &common_msg.SendProtocol{
			Protocol: &resp.MoveMob{
				IsAggroed:   m.IsAggroed,
				CenterSplit: m.CenterSplit,
				Skill1:      m.Skill1,
				Skill2:      m.Skill2,
				Skill3:      m.Skill3,
				Skill4:      m.Skill4,
				OID:         state.OID,
				StartPoint:  state.Position,
				Movements:   m.Movements,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		},
	})
}
