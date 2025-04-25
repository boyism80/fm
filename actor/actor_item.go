package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

type ItemActor struct {
	Item    *model.Item
	handler *handler.MessageHandler
}

func NewItemActor(Item *model.Item) actor.Actor {
	act := &ItemActor{
		Item:    Item,
		handler: handler.NewMessageHandler(),
	}
	handler.RegisterItemHandlers(act.Item, act.handler)
	return act
}

func (state *ItemActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
