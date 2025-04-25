package model

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type MainActor struct {
	handler *handler.MessageHandler
}

func NewMainActor() actor.Actor {
	act := &MainActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterMainHandlers(act, act.handler)
	return act
}

func (state *MainActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
