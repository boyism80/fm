package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
)

type CharacterLootFailed struct {
	Oid uint32
}

type CharacterItemLooting struct {
	Pid  *actor.PID
	Oid  uint32
	Item entity.Item
}

type CharacterMesoLooting struct {
	Pid  *actor.PID
	Oid  uint32
	Meso int32
}
