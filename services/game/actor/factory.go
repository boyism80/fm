package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type Factory struct {
	gw game.GameWorld
}

func NewFactory(gw game.GameWorld) *Factory {
	return &Factory{gw: gw}
}

func (f *Factory) Start() error { return nil }
func (f *Factory) Stop()        {}

func (f *Factory) NewMapActor(m *game.Map) actor.Actor {
	return NewMapActor(m, f.gw)
}

func (f *Factory) NewStateMachineActor(sm *game.StateMachine) actor.Actor {
	return NewStateMachineActor(sm, f.gw)
}
