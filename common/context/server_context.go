package context

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/data"
)

type ServerContext struct {
	actor.Context
	Resources *data.Resources
	MapActors map[uint32]*actor.PID
}

func NewServerContext(
	Resources *data.Resources,
	mapActors map[uint32]*actor.PID,
	serverCtxActor *actor.PID) *ServerContext {

	return &ServerContext{
		Resources: Resources,
		MapActors: mapActors,
	}
}
