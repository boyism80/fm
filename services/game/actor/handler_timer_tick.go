package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"log"
)

type TimerTickHandler struct{}

func (TimerTickHandler) New() *TimerTickHandler {
	return &TimerTickHandler{}
}

func (h *TimerTickHandler) Handle(ctx actor.Context, a *MapActor, msg *TimerTick) {
	handler := a.timerReg.GetHandler(msg.HandlerName)
	if handler == nil {
		return
	}
	for _, m := range a.Maps() {
		if m == nil {
			continue
		}
		pid := m.GetActorPID()
		if pid == nil || !pid.Equal(ctx.Self()) {
			continue
		}
		if err := handler.Handle(ctx, m); err != nil {
			log.Printf("Timer handler %s error: %v", handler.GetName(), err)
		}
	}
}
