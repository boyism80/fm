package actorx

import (
	"math"
	"math/rand/v2"

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
	handler   *handler.MessageHandler
	mapPID    *actor.PID
	ServerCtx *context.ServerContext
}

func NewMobActor(ctx actor.Context, serverCtx *context.ServerContext, entity entity.Mob, mapPID *actor.PID) actor.Actor {
	actor := &MobActor{
		Mob:       entity,
		handler:   handler.NewMessageHandler(),
		mapPID:    mapPID,
		ServerCtx: serverCtx,
	}
	RegisterLifeHandlers(ctx, &actor.Mob.Life, actor.handler)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobSpawn)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobPlayerWarped)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobKill)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobControllerChange)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobMove)
	handler.RegisterHandler(ctx, actor, actor.handler, onMobDamaged)
	return actor
}

func (state *MobActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}

func onMobSpawn(ctx actor.Context, state *MobActor, m *msg.Spawn) {
	ctx.Send(state.mapPID, &msg.MapBroadcast{
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

	ctx.Send(state.mapPID, &msg.MapSpawnedMob{
		PID: ctx.Self(),
		OID: state.OID,
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

	if m.After != nil {
		ctx.Send(m.After, &common_msg.SendProtocol{
			Protocol: &resp.StartControlMob{
				Mob:   &state.Mob,
				Aggro: false,
			},
			Policy: types.SEND_POLICY_ENCRYPT,
		})
	} else {
		ctx.Send(m.Before, &common_msg.SendProtocol{
			Protocol: &resp.StopControlMob{
				OID: state.OID,
			},
		})
	}
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

	ctx.Send(state.mapPID, &msg.MapBroadcast{
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

func onMobDamaged(ctx actor.Context, state *MobActor, m *msg.MobDamaged) {

	if state.Hp <= 0 {
		return
	}

	isDead := false
	for _, damage := range m.DamagePairs {
		state.Hp -= uint16(math.Min(float64(state.Hp), float64(damage.Damage)))
		if state.Hp == 0 {
			isDead = true
			break
		}
	}

	ctx.Send(m.Sender, &common_msg.SendProtocol{
		Protocol: &resp.ShowMobHp{
			OID:        state.OID,
			Percentage: uint8(state.Hp * 100 / state.MaxHp),
		},
		Policy: types.SEND_POLICY_ENCRYPT,
	})

	// TODO: drop items

	if isDead {

		drops, ok := state.ServerCtx.Resources.Drops[state.Spec.ID]
		if ok {
			for _, drop := range drops {
				// drop.Prob 확률로 아이템 드랍
				if rand.Float32() > drop.Prob {
					continue
				}

				if drop.Item == 0 {
					min := float64(drop.Money) * 0.75
					max := float64(drop.Money)
					count := int32(min + rand.Float64()*(max-min))
					if count == 0 {
						continue
					}

					ctx.Send(state.mapPID, &msg.MapSpawnMeso{
						Count:        count,
						SpawnedPoint: state.Position,
						Owner:        m.Sender,
						OwnerID:      state.OID,
					})
				} else {
					count := uint16(1)
					if drop.Max != 0 && drop.Min != 0 {
						count = uint16(rand.Uint32N(uint32(drop.Max-drop.Min)+1) + uint32(drop.Min))
					}

					item, err := entity.NewItem(state.ServerCtx, drop.Item, count)
					if err != nil {
						continue
					}
					item.BindDrop(&entity.Drop{
						Object:       &entity.Object{},
						Owner:        m.CharacterId,
						SpawnedPoint: state.Position,
						DropType:     constant.DropTypeOwned,
						Looting:      false,
					})
					ctx.Send(state.mapPID, &msg.MapSpawnItem{
						Item:    item,
						Owner:   m.Sender,
						OwnerID: state.OID,
					})
				}
			}
		}

		ctx.Send(m.Sender, &msg.CharacterKillMob{
			OID:   state.OID,
			MobID: state.Spec.ID,
		})

		ctx.Send(state.mapPID, &msg.MapDieMob{
			Sender:        ctx.Self(),
			OID:           state.OID,
			AnimationType: constant.MobDieAnimationTypeFadeOut,
			Position:      state.Position,
		})
	}
}
