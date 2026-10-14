package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
	"github.com/boyism80/fm/core/route"
	"github.com/boyism80/fm/services/game/actor/timers"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/domain"
)

type GameLogicActor struct {
	GameWorld    domain.GameWorld
	scheduler    *scheduler.TimerScheduler
	timerReg     *TimerRegistry
	timerCancels []scheduler.CancelFunc
	maps         func() []*domain.Map
}

func (a *GameLogicActor) Receive(ctx actor.Context) {
	msg := ctx.Message()
	if env, ok := msg.(*route.Envelope); ok {
		a.handleEnvelope(ctx, env)
		return
	}
	mapMessageRegistry.Dispatch(ctx, a, msg)
}

func (a *GameLogicActor) Maps() []*domain.Map {
	if a == nil || a.maps == nil {
		return nil
	}
	return a.maps()
}

func (a *GameLogicActor) GetCharacter(characterID uint32) *domain.Map {
	for _, m := range a.Maps() {
		if m != nil && m.GetPlayer(characterID) != nil {
			return m
		}
	}
	return nil
}

func (a *GameLogicActor) GetObject(objectType constant.ObjectType, oid uint32) *domain.Map {
	for _, m := range a.Maps() {
		if m != nil && m.GetObject(objectType, oid) != nil {
			return m
		}
	}
	return nil
}

func (a *GameLogicActor) handleEnvelope(ctx actor.Context, env *route.Envelope) {
	if env == nil {
		return
	}
	inner := env.Inner
	if inner == nil {
		a.finishEnvelope(ctx, env, false, "nil_inner")
		return
	}
	if !a.hasCharacterOnMap(env.CharacterID) {
		if a.GameWorld != nil {
			a.GameWorld.Systems().Get[*domain.DispatchSystem]().Resend(env)
		}
		return
	}
	mapMessageRegistry.Dispatch(ctx, a, inner)
	a.finishEnvelope(ctx, env, true, "")
}

func (a *GameLogicActor) hasCharacterOnMap(characterID uint32) bool {
	if a == nil || characterID == 0 {
		return false
	}
	return a.GetCharacter(characterID) != nil
}

func (a *GameLogicActor) registerTimers() {
	a.timerReg.Register(&timers.MobSpawnTimer{})
	a.timerReg.Register(&timers.ItemCleanupTimer{})
	a.timerReg.Register(&timers.CooldownCheckTimer{})
	a.timerReg.Register(&timers.ItemExpireTimer{})
	a.timerReg.Register(&timers.ClientPingTimer{})
	a.timerReg.Register(&timers.MobBuffExpireTimer{})
	a.timerReg.Register(&timers.MobPoisonTickTimer{})
	a.timerReg.Register(&timers.MistExpireTimer{})
	a.timerReg.Register(&timers.MistPoisonTickTimer{})
	a.timerReg.Register(&timers.CharacterSaveTimer{})
	a.timerReg.Register(&timers.PartySearchTimer{})
}

func (a *GameLogicActor) StopTimers() {
	for _, cancel := range a.timerCancels {
		if cancel != nil {
			cancel()
		}
	}
	a.timerCancels = nil
}

func (a *GameLogicActor) finishEnvelope(ctx actor.Context, msg *route.Envelope, ok bool, reason string) {
	if msg == nil {
		return
	}
	if a.GameWorld != nil {
		a.GameWorld.Systems().Get[*domain.DispatchSystem]().CompleteEnvelope(msg.CorrelationID)
	}
	if msg.Caller == nil {
		return
	}
	ctx.Send(msg.Caller, &route.Result{
		CorrelationID: msg.CorrelationID,
		OK:            ok,
		Reason:        reason,
	})
}
