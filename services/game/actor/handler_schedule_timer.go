package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"time"
)

type ScheduleTimerHandler struct{}

func (ScheduleTimerHandler) New() *ScheduleTimerHandler {
	return &ScheduleTimerHandler{}
}

func (h *ScheduleTimerHandler) Handle(ctx actor.Context, a *MapActor, msg *c_actor.ScheduleTimer) {
	if msg.Logic == nil {
		return
	}

	go func() {
		time.Sleep(msg.Interval)
		ctx.Send(ctx.Self(), &c_actor.ExecuteTimer{
			Logic: msg.Logic,
		})
	}()
}
