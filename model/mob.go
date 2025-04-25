package model

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

type MobActor struct {
	Life
	Aggro   bool
	handler *handler.MessageHandler
}

func NewMobActor() actor.Actor {
	act := &MobActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterMobHandlers(act, act.handler)
	return act
}

func (state *MobActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
