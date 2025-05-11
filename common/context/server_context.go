package context

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/data"
)

type ServerContext struct {
	Resources *data.Resources
	MapActors map[uint32]*protoactor.PID
}

func NewServerContext(
	Resources *data.Resources,
	mapActors map[uint32]*protoactor.PID,
	serverCtxActor *protoactor.PID) *ServerContext {

	return &ServerContext{
		Resources: Resources,
		MapActors: mapActors,
	}
}
