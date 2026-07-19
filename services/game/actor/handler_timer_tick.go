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
	if a.Map == nil {
		return
	}
	handler := a.timerReg.GetHandler(msg.HandlerName)
	if handler == nil {
		return
	}
	if err := handler.Handle(ctx, a.Map); err != nil {
		log.Printf("Timer handler %s error: %v", handler.GetName(), err)
	}
}
