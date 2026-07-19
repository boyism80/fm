package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core/ensure"
	"github.com/boyism80/fm/services/game/actor/timers"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type GameLogicActor struct {
	GameWorld    entity.GameWorld
	scheduler    *scheduler.TimerScheduler
	timerReg     *TimerRegistry
	timerCancels []scheduler.CancelFunc
	maps         func() []*entity.Map
}

func (a *GameLogicActor) Receive(ctx actor.Context) {
	msg := ctx.Message()
	if env, ok := msg.(*ensure.EnsureDeliver); ok {
		a.handleEnsureDeliver(ctx, env)
		return
	}
	mapMessageRegistry.Dispatch(ctx, a, msg)
}

func (a *GameLogicActor) Maps() []*entity.Map {
	if a == nil || a.maps == nil {
		return nil
	}
	return a.maps()
}

func (a *GameLogicActor) GetCharacter(characterID uint32) *entity.Map {
	for _, m := range a.Maps() {
		if m != nil && m.GetPlayer(characterID) != nil {
			return m
		}
	}
	return nil
}

func (a *GameLogicActor) GetObject(objectType constant.ObjectType, oid uint32) *entity.Map {
	for _, m := range a.Maps() {
		if m != nil && m.GetObject(objectType, oid) != nil {
			return m
		}
	}
	return nil
}

func (a *GameLogicActor) handleEnsureDeliver(ctx actor.Context, env *ensure.EnsureDeliver) {
	if env == nil {
		return
	}
	inner := env.Inner
	if inner == nil {
		a.ensureFinish(ctx, env, false, "nil_inner")
		return
	}
	if !a.hasCharacterOnMap(env.CharacterID) {
		if a.GameWorld != nil {
			a.GameWorld.EnsureRedispatch(env)
		}
		return
	}
	mapMessageRegistry.Dispatch(ctx, a, inner)
	a.ensureFinish(ctx, env, true, "")
}

func (a *GameLogicActor) hasCharacterOnMap(characterID uint32) bool {
	if a == nil || characterID == 0 {
		return false
	}
	return a.GetCharacter(characterID) != nil
}

func (a *GameLogicActor) registerTimers() {
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

func (a *GameLogicActor) StopTimers() {
	for _, cancel := range a.timerCancels {
		if cancel != nil {
			cancel()
		}
	}
	a.timerCancels = nil
}

func (a *GameLogicActor) ensureFinish(ctx actor.Context, msg *ensure.EnsureDeliver, ok bool, reason string) {
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
