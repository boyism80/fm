package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
)

type CharacterLootFailed struct {
	OID uint32
}

type CharacterItemLooting struct {
	PID  *actor.PID
	OID  uint32
	Item entity.Item
}

type CharacterMesoLooting struct {
	PID  *actor.PID
	OID  uint32
	Meso int32
}

type CharacterMapChanged struct {
	MID        uint32
	Map        *actor.PID
	Init       bool
	SpawnPoint uint8
}

type RunScript struct {
	Script string
}
