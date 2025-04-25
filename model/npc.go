package model

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type NpcActor struct {
	Object
	Dialog  []string
	handler *handler.MessageHandler
}

func NewNpcActor() actor.Actor {
	act := &NpcActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterNpcHandlers(act, act.handler)
	return act
}

func (state *NpcActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
