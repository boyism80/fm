package model

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type ItemActor struct {
	Object
	ID      int64
	Name    string
	handler *handler.MessageHandler
}

func NewItemActor() actor.Actor {
	act := &ItemActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterItemHandlers(act, act.handler)
	return act
}

func (state *ItemActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
