package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type MapActor struct {
	GameLogicActor
	Map *entity.Map
}

func NewMapActor(m *entity.Map, gameWorld entity.GameWorld) *MapActor {
	a := &MapActor{Map: m}
	a.GameWorld = gameWorld
	a.maps = func() []*entity.Map {
		if a.Map == nil {
			return nil
		}
		return []*entity.Map{a.Map}
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
