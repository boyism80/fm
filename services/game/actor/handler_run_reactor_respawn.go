package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
)

type RunReactorRespawnHandler struct{}

func (RunReactorRespawnHandler) New() *RunReactorRespawnHandler {
	return &RunReactorRespawnHandler{}
}

func (h *RunReactorRespawnHandler) Handle(ctx actor.Context, a *MapActor, msg *c_actor.RunReactorRespawn) {
	if a.Map == nil || msg == nil {
		return
	}

	reactorSpawn := a.Map.GetReactorSpawn(msg.SpawnID)
	if reactorSpawn == nil || reactorSpawn.Spawned {
		return
	}

	_, _ = a.Map.SpawnReactor(reactorSpawn)
}
