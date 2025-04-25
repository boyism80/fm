package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

type MainActor struct {
	Main    *model.Main
	handler *handler.MessageHandler
}

func NewMainActor() protoactor.Actor {
	act := &MainActor{
		Main:    &model.Main{},
		handler: handler.NewMessageHandler(),
	}
	RegisterMainHandlers(act.Main, act.handler)
	return act
}

func (state *MainActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
