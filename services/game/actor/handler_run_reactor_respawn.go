package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type RunReactorRespawnHandler struct{}

func (RunReactorRespawnHandler) New() *RunReactorRespawnHandler {
	return &RunReactorRespawnHandler{}
}

func (h *RunReactorRespawnHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *c_actor.RunReactorRespawn) {
	if msg == nil {
		return
	}
	var targetMap *entity.Map
	for _, m := range a.Maps() {
		if m != nil && m.GetMapID() == msg.MapID {
			targetMap = m
			break
		}
	}
	if targetMap == nil {
		return
	}
	pid := targetMap.LogicActorPID()
	if pid == nil || !pid.Equal(ctx.Self()) {
		return
	}
	reactorSpawn := targetMap.GetReactorSpawn(msg.SpawnID)
	if reactorSpawn == nil || reactorSpawn.Spawned {
		return
	}
	_, _ = targetMap.SpawnReactor(reactorSpawn)
}
