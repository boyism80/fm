package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
)

type ItemActor struct {
	entity.Item
	handler *handler.MessageHandler
}

func NewItemActor() protoactor.Actor {
	act := &ItemActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterItemHandlers(act, act.handler)
	return act
}

func (state *ItemActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
