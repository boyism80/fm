package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
)

type ItemSpawn struct {
	OwnerId      uint32
	SpawnedPoint types.Vector2[int16]
}

type ItemLooting struct {
	Actor    *actor.PID
	Position types.Vector2[int16]
}

type ItemLooted struct {
	Success     bool
	CharacterId uint32
	Count       int32
}
