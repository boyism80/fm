package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

type ObjectActor struct {
	Object  *model.Object
	handler *handler.MessageHandler
}

func NewObjectActor(object *model.Object) actor.Actor {
	act := &ObjectActor{
		Object:  object,
		handler: handler.NewMessageHandler(),
	}
	handler.RegisterObjectHandlers(act.Object, act.handler)
	return act
}

func (state *ObjectActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
