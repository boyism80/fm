package entity

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/game/wz"
)

type GameContext interface {
	GetResources() *wz.Resources
	GetMap(mapId uint32) *Map
	GetExpRate() int
	GetDropRate() int
	GetMesoRate() int
	RequestWarp(character *Character, targetMap *Map, spawnPoint uint8) error
	DispatchRunCharacterTimer(pid *actor.PID, payload *c_actor.RunCharacterTimer)
	NotifyDoorSpawn(returnMapWZID uint32, spawn DoorSpawn)
	NotifyDoorRemove(removal DoorRemove)
}
