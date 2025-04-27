package model

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/dao"
	"github.com/boyism80/fm/handler"
)

type ItemActor struct {
	dao.Item
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
