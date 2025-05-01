package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type ServerActor struct {
	handler *handler.MessageHandler
}

func NewMainActor() protoactor.Actor {
	act := &ServerActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterMainHandlers(nil, act, act.handler)
	return act
}

func (state *ServerActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
