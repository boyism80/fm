package actorx

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
)

type LoginServerActor struct {
	handler *handler.MessageHandler
}

func NewLoginServerActor() actor.Actor {
	actor := &LoginServerActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterLoginServerHandlers(nil, actor, actor.handler)
	return actor
}

func (state *LoginServerActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
