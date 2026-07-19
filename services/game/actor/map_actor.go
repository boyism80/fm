package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core/ensure"
	"github.com/boyism80/fm/services/game/actor/timers"
	"github.com/boyism80/fm/services/game/entity"
)

type MapActor struct {
	Map       *entity.Map
	GameWorld entity.GameWorld
	scheduler *scheduler.TimerScheduler
	timerReg  *TimerRegistry
}

func (a *MapActor) Receive(ctx actor.Context) {
	msg := ctx.Message()
	if env, ok := msg.(*ensure.EnsureDeliver); ok {
		a.handleEnsureDeliver(ctx, env)
		return
	}
	a.dispatch(ctx, msg)
}

func (a *MapActor) dispatch(ctx actor.Context, msg interface{}) {
	mapMessageRegistry.Dispatch(ctx, a, msg)
}

func (a *MapActor) handleEnsureDeliver(ctx actor.Context, env *ensure.EnsureDeliver) {
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
	a.dispatch(ctx, inner)
	a.ensureFinish(ctx, env, true, "")
}

func (a *MapActor) hasCharacterOnMap(characterID uint32) bool {
	if a.Map == nil || characterID == 0 {
		return false
	}
	return a.Map.GetPlayer(characterID) != nil
}

func (a *MapActor) registerTimers() {
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

func (a *MapActor) ensureNotOnMap(msg *ensure.EnsureDeliver) {
	if a.GameWorld != nil {
		a.GameWorld.EnsureRedispatch(msg)
	}
}

func (a *MapActor) ensureFinish(ctx actor.Context, msg *ensure.EnsureDeliver, ok bool, reason string) {
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
