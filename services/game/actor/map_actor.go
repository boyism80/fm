package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/ensure"
	"github.com/boyism80/fm/services/game/actor/timers"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type LogicActor struct {
	Map          *entity.Map
	GameWorld    entity.GameWorld
	scheduler    *scheduler.TimerScheduler
	timerReg     *TimerRegistry
	StateMachine *entity.StateMachine
}

type MapActor = LogicActor

func (a *LogicActor) Owner() *LogicActor {
	return a
}

func (a *LogicActor) Receive(ctx actor.Context) {
	msg := ctx.Message()
	if target := a.forwardTarget(); target != nil && a.shouldForward(msg) {
		ctx.Send(target, msg)
		return
	}
	if env, ok := msg.(*ensure.EnsureDeliver); ok {
		a.handleEnsureDeliver(ctx, env)
		return
	}
	a.dispatch(ctx, msg)
}

func (a *LogicActor) forwardTarget() *actor.PID {
	if a == nil || a.StateMachine != nil || a.Map == nil {
		return nil
	}
	sm := a.Map.StateMachine()
	if sm == nil {
		return nil
	}
	return sm.ActorPID
}

func (a *LogicActor) shouldForward(msg interface{}) bool {
	switch msg.(type) {
	case *actor.Started, *actor.Stopping, *actor.Stopped, *actor.Restarting:
		return false
	case *AttachStateMachine, *DetachStateMachine:
		return false
	case *TimerTick:
		return false
	}
	return true
}

func (a *LogicActor) dispatch(ctx actor.Context, msg interface{}) {
	a.selectMap(msg)
	mapMessageRegistry.Dispatch(ctx, a, msg)
}

func (a *LogicActor) Maps() []*entity.Map {
	if a == nil {
		return nil
	}
	if a.StateMachine != nil {
		return a.StateMachine.MapList()
	}
	if a.Map == nil {
		return nil
	}
	return []*entity.Map{a.Map}
}

func (a *LogicActor) MapForCharacter(characterID uint32) *entity.Map {
	for _, m := range a.Maps() {
		if m != nil && m.GetPlayer(characterID) != nil {
			return m
		}
	}
	return nil
}

func (a *LogicActor) MapForObject(objectType constant.ObjectType, oid uint32) *entity.Map {
	for _, m := range a.Maps() {
		if m != nil && m.GetObject(objectType, oid) != nil {
			return m
		}
	}
	return nil
}

func (a *LogicActor) selectMap(msg interface{}) {
	if a == nil || a.StateMachine == nil {
		return
	}
	switch v := msg.(type) {
	case *WarpCharacter:
		a.Map = v.TargetMap
	case *AddCharacter:
		a.Map = v.TargetMap
	case *HandoffCharacter:
		if v.Character == nil {
			a.Map = nil
		} else {
			a.Map = a.MapForCharacter(v.Character.GetID())
		}
	case *RemoveCharacter:
		a.Map = a.MapForCharacter(v.CharacterID)
	case *PartyMemberLeft:
		a.Map = a.MapForCharacter(v.LeaverID)
	case *ResponseSpawnDoor:
		a.Map = a.MapForCharacter(v.CharacterID)
	case *c_actor.RunObjectTimer:
		a.Map = a.MapForObject(constant.ObjectType(v.ObjectType), v.ID)
	default:
		a.Map = nil
	}
}

func (a *LogicActor) handleEnsureDeliver(ctx actor.Context, env *ensure.EnsureDeliver) {
	if env == nil {
		return
	}
	inner := env.Inner
	if inner == nil {
		a.ensureFinish(ctx, env, false, "nil_inner")
		return
	}
	if !a.hasCharacterOnMap(env.CharacterID) {
		a.ensureNotOnMap(env)
		return
	}
	if a.StateMachine != nil {
		a.Map = a.MapForCharacter(env.CharacterID)
	}
	mapMessageRegistry.Dispatch(ctx, a, inner)
	a.ensureFinish(ctx, env, true, "")
}

func (a *LogicActor) hasCharacterOnMap(characterID uint32) bool {
	if a == nil || characterID == 0 {
		return false
	}
	return a.MapForCharacter(characterID) != nil
}

func (a *LogicActor) registerTimers() {
	RegisterTimer[*timers.MobSpawnTimer](a.timerReg)
	RegisterTimer[*timers.ItemCleanupTimer](a.timerReg)
	RegisterTimer[*timers.CooldownCheckTimer](a.timerReg)
	RegisterTimer[*timers.ClientPingTimer](a.timerReg)
	RegisterTimer[*timers.BuffExpireTimer](a.timerReg)
	RegisterTimer[*timers.MobBuffExpireTimer](a.timerReg)
	RegisterTimer[*timers.MobPoisonTickTimer](a.timerReg)
	RegisterTimer[*timers.MistExpireTimer](a.timerReg)
	RegisterTimer[*timers.MistPoisonTickTimer](a.timerReg)
	RegisterTimer[*timers.CharacterSaveTimer](a.timerReg)
	RegisterTimer[*timers.PartySearchTimer](a.timerReg)
}

func (a *LogicActor) ensureNotOnMap(msg *ensure.EnsureDeliver) {
	if a.GameWorld != nil {
		a.GameWorld.EnsureRedispatch(msg)
	}
}

func (a *LogicActor) ensureFinish(ctx actor.Context, msg *ensure.EnsureDeliver, ok bool, reason string) {
	if msg == nil {
		return
	}
	if a.GameWorld != nil {
		a.GameWorld.EnsureComplete(msg.CorrelationID)
	}
	if msg.Caller == nil {
		return
	}
	ctx.Send(msg.Caller, &ensure.EnsureResult{
		CorrelationID: msg.CorrelationID,
		OK:            ok,
		Reason:        reason,
	})
}
