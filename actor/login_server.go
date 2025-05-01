package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type LoginServerActor struct {
	handler *handler.MessageHandler
}

func NewLoginServerActor() protoactor.Actor {
	act := &LoginServerActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterLoginServerHandlers(nil, act, act.handler)
	return act
}

func (state *LoginServerActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
