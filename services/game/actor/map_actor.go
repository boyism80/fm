package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type MapActor struct {
	GameLogicActor
	Map *game.Map
}

func NewMapActor(m *game.Map, gameWorld game.GameWorld) *MapActor {
	a := &MapActor{Map: m}
	a.GameWorld = gameWorld
	a.maps = func() []*game.Map {
		if a.Map == nil {
			return nil
		}
		return []*game.Map{a.Map}
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
			a.Map.ClearObjectTimers()
			a.Map.ClearLuaRoot()
		}
	default:
		a.GameLogicActor.Receive(ctx)
	}
}
