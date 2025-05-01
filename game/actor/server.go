package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
)

type GameServerActor struct {
	handler *handler.MessageHandler
}

func NewGameServerActor() protoactor.Actor {
	act := &GameServerActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterGameServerHandlers(nil, act, act.handler)
	return act
}

func (state *GameServerActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
