package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
)

type ResumeLuaHandler struct{}

func (ResumeLuaHandler) New() *ResumeLuaHandler {
	return &ResumeLuaHandler{}
}

func (h *ResumeLuaHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *ResumeLua) {
	if msg.Root == nil || msg.Thread == nil {
		return
	}
	resumeArgs := make([]interface{}, len(msg.Args))
	for i, a := range msg.Args {
		resumeArgs[i] = a
	}
	luax.Resume(msg.Root, msg.Thread, resumeArgs...)
}
