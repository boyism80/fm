package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type Factory struct {
	gw domain.GameWorld
}

func NewFactory(gw domain.GameWorld) *Factory {
	return &Factory{gw: gw}
}

func (f *Factory) Start() error { return nil }
func (f *Factory) Stop()        {}

func (f *Factory) NewMapActor(m *domain.Map) actor.Actor {
	return NewMapActor(m, f.gw)
}

func (f *Factory) NewStateMachineActor(sm *domain.StateMachine) actor.Actor {
	return NewStateMachineActor(sm, f.gw)
}
