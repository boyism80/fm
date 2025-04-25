package model

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type MainActor struct {
	handler *handler.MessageHandler
}

func NewMainActor() protoactor.Actor {
	act := &MainActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterMainHandlers(act, act.handler)
	return act
}

func (state *MainActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
