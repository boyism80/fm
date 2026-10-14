package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type MapActor struct {
	GameLogicActor
	Map *domain.Map
}

func NewMapActor(m *domain.Map, gameWorld domain.GameWorld) *MapActor {
	a := &MapActor{Map: m}
	a.GameWorld = gameWorld
	a.maps = func() []*domain.Map {
		if a.Map == nil {
			return nil
		}
		return []*domain.Map{a.Map}
	}
	return a
}

func (a *MapActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case *actor.Restarting:
		a.StopTimers()
	case *actor.Stopped:
		a.StopTimers()
		if a.Map != nil {
			a.Map.ClearLuaRoot()
		}
	default:
		a.GameLogicActor.Receive(ctx)
	}
}
