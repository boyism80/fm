package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type RunObjectTimerHandler struct{}

func (RunObjectTimerHandler) New() *RunObjectTimerHandler {
	return &RunObjectTimerHandler{}
}

func (h *RunObjectTimerHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *c_actor.RunObjectTimer) {
	if msg == nil {
		return
	}
	m := a.GetObject(constant.ObjectType(msg.ObjectType), msg.ID)
	if m == nil {
		return
	}
	pid := m.GetActorPID()
	if pid == nil || !pid.Equal(ctx.Self()) {
		return
	}
	obj := m.GetObject(constant.ObjectType(msg.ObjectType), msg.ID)
	if obj == nil {
		return
	}
	entry := obj.GetTimerEntry(msg.Key)
	if entry == nil {
		return
	}
	if entry.Callback != nil {
		entry.Callback()
	}
	if entry.Repeat {
		obj.RescheduleTimer(msg.Key)
	} else {
		obj.RemoveTimer(msg.Key)
	}
}
