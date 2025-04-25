package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

type ClientActor struct {
	Client  *model.Client
	handler *handler.MessageHandler
}

func NewClientActor(Client *model.Client) actor.Actor {
	act := &ClientActor{
		Client:  Client,
		handler: handler.NewMessageHandler(),
	}
	return act
}

func (state *ClientActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
