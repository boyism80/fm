package actor

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
)

type TimerHandler interface {
	GetName() string
	GetInterval() time.Duration
	GetInitialDelay() time.Duration
	Handle(ctx actor.Context, mapData *entity.Map) error
}

type TimerRegistry struct {
	handlers map[string]TimerHandler
}

func NewTimerRegistry() *TimerRegistry {
	return &TimerRegistry{
		handlers: make(map[string]TimerHandler),
	}
}

func RegisterTimer[H interface {
	TimerHandler
	New() H
}](registry *TimerRegistry) {
	var zero H
	handler := zero.New()
	registry.handlers[handler.GetName()] = handler
}

func (r *TimerRegistry) GetAllHandlers() []TimerHandler {
	out := make([]TimerHandler, 0, len(r.handlers))
	for _, h := range r.handlers {
		out = append(out, h)
	}
	return out
}

func (r *TimerRegistry) GetHandler(name string) TimerHandler {
	return r.handlers[name]
}
