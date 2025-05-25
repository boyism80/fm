package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/types"
)

type CharacterMove struct {
	Position types.Vector2[int16]
}

type Warped struct {
	Sender *actor.PID
}
