package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
)

type Spawn struct {
}

type ItemLooting struct {
	Actor       *actor.PID
	Position    types.Vector2[int16]
	CharacterId uint32
}

type ItemLooted struct {
	Success     bool
	CharacterId uint32
	Count       int32
}

type ItemDropTypeChanged struct {
	Mode constant.DropType
}

type ItemDestroy struct {
	Animation constant.DropItemAnimationType
}
