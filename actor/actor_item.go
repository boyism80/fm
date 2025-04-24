package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/model"
)

type ItemActor struct {
	Item    *model.Item
	Context actor.Context
}
