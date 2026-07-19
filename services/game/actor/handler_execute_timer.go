package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"log"
)

type ExecuteTimerHandler struct{}

func (ExecuteTimerHandler) New() *ExecuteTimerHandler {
	return &ExecuteTimerHandler{}
}

func (h *ExecuteTimerHandler) Handle(ctx actor.Context, a *MapActor, msg *c_actor.ExecuteTimer) {
	if msg.Logic == nil {
		return
	}

	if err := msg.Logic(); err != nil {
		log.Printf("Error executing timer logic: %v", err)
	}
}
