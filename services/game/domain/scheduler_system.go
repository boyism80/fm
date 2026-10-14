package domain

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
)

type SchedulerSystem struct {
	gw GameWorld
}

func NewSchedulerSystem(gw GameWorld) *SchedulerSystem {
	return &SchedulerSystem{gw: gw}
}

func (s *SchedulerSystem) Start() error { return nil }
func (s *SchedulerSystem) Stop()        {}

func (s *SchedulerSystem) RunObjectTimer(pid *actor.PID, obj Object, key string) {
	if s.gw == nil || pid == nil || obj == nil || key == "" {
		return
	}
	m := obj.GetMap()
	if m == nil {
		return
	}
	payload := &c_actor.RunObjectTimer{
		MapID:      m.GetMapID(),
		ObjectType: obj.GetObjectType(),
		ID:         obj.GetPK(),
		Key:        key,
	}
	if root := s.gw.Systems().Get[*c_actor.ActorRegistry]().Root(); root != nil {
		root.Send(pid, payload)
	}
}

func (s *SchedulerSystem) RunReactorRespawn(pid *actor.PID, mapID uint32, spawnID uint32) {
	if s.gw == nil || pid == nil {
		return
	}
	payload := &c_actor.RunReactorRespawn{
		MapID:   mapID,
		SpawnID: spawnID,
	}
	if root := s.gw.Systems().Get[*c_actor.ActorRegistry]().Root(); root != nil {
		root.Send(pid, payload)
	}
}
