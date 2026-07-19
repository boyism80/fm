package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

type ResumeLuaHandler struct{}

func (ResumeLuaHandler) New() *ResumeLuaHandler {
	return &ResumeLuaHandler{}
}

func (h *ResumeLuaHandler) Handle(ctx actor.Context, a *MapActor, msg *ResumeLua) {
	if msg.Root == nil || msg.Thread == nil {
		return
	}
	resumeArgs := make([]interface{}, len(msg.Args))
	for i, a := range msg.Args {
		resumeArgs[i] = a
	}
	state, _, _ := luax.Resume(msg.Root, msg.Thread, "", resumeArgs...)
	if state == lua.ResumeOK {
		cfg, ok := luax.GetConfiguration(msg.Thread)
		if ok && cfg.KeepAlive {
			luax.Close(msg.Thread)
		}
	}
}
