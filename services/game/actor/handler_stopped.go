package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type StoppedHandler struct{}

func (StoppedHandler) New() *StoppedHandler {
	return &StoppedHandler{}
}

func (h *StoppedHandler) Handle(ctx actor.Context, a *MapActor, _ *actor.Stopped) {
	if a.StateMachine == nil && a.Map != nil {
		a.Map.ClearLuaRoot()
	}
}
