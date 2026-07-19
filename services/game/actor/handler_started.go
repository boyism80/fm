package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/scheduler"
)

type StartedHandler struct{}

func (StartedHandler) New() *StartedHandler {
	return &StartedHandler{}
}

func (h *StartedHandler) Handle(ctx actor.Context, a *GameLogicActor, _ *actor.Started) {
	a.scheduler = scheduler.NewTimerScheduler(ctx)
	a.timerReg = NewTimerRegistry()
	a.registerTimers()

	for _, handler := range a.timerReg.GetAllHandlers() {
		a.scheduler.SendRepeatedly(
			handler.GetInitialDelay(),
			handler.GetInterval(),
			ctx.Self(),
			&TimerTick{
				HandlerName: handler.GetName(),
			},
		)
	}
}
