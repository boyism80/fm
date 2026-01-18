package actor

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
)

// TimerHandler defines the interface for timer handlers
type TimerHandler interface {
	GetName() string
	GetInterval() time.Duration
	GetInitialDelay() time.Duration
	Handle(ctx actor.Context, mapData *entity.Map) error
}

// TimerHandlerConstructor defines the factory interface for timer handlers
type TimerHandlerConstructor[H TimerHandler] interface {
	New() H
}

// TimerRegistry manages timer handlers for MapActor
type TimerRegistry struct {
	handlers []TimerHandler
}

func NewTimerRegistry() *TimerRegistry {
	return &TimerRegistry{
		handlers: make([]TimerHandler, 0),
	}
}

// RegisterTimer registers a timer handler (similar to core.Bind pattern)
func RegisterTimer[H TimerHandler, C TimerHandlerConstructor[H]](registry *TimerRegistry) {
	var constructor C
	handler := constructor.New()
	registry.handlers = append(registry.handlers, handler)
}

// GetAllHandlers returns all registered timer handlers
func (r *TimerRegistry) GetAllHandlers() []TimerHandler {
	return r.handlers
}
