package actor

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type TimerTickHandler struct{}

func (TimerTickHandler) New() *TimerTickHandler {
	return &TimerTickHandler{}
}

func (h *TimerTickHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *TimerTick) {
	handler := a.timerReg.GetHandler(msg.HandlerName)
	if handler == nil {
		return
	}
	for _, m := range a.Maps() {
		if m == nil {
			continue
		}
		if err := handler.Handle(ctx, m); err != nil {
			log.Printf("Timer handler %s error: %v", handler.GetName(), err)
		}
	}
}
