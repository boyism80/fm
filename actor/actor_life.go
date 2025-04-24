package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

type LifeActor struct {
	Life    *model.Life
	handler *handler.MessageHandler
}

func NewLifeActor(life *model.Life) actor.Actor {
	act := &LifeActor{
		Life:    life,
		handler: handler.NewMessageHandler(),
	}
	handler.RegisterLifeHandlers(act.Life, act.handler)
	return act
}

func (state *LifeActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
