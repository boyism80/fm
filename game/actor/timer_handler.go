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

// TimerRegistry manages timer handlers for MapActor
type TimerRegistry struct {
	handlers []TimerHandler
}

func NewTimerRegistry() *TimerRegistry {
	return &TimerRegistry{
		handlers: make([]TimerHandler, 0),
	}
}

// RegisterTimer registers a timer handler. H must implement TimerHandler and New() H.
func RegisterTimer[H interface {
	TimerHandler
	New() H
}](registry *TimerRegistry) {
	var zero H
	registry.handlers = append(registry.handlers, zero.New())
}

// GetAllHandlers returns all registered timer handlers
func (r *TimerRegistry) GetAllHandlers() []TimerHandler {
	return r.handlers
}
