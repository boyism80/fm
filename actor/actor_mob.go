package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

type MobActor struct {
	Mob     *model.Mob
	handler *handler.MessageHandler
}

func NewMobActor(Mob *model.Mob) actor.Actor {
	act := &MobActor{
		Mob:     Mob,
		handler: handler.NewMessageHandler(),
	}
	handler.RegisterMobHandlers(act.Mob, act.handler)
	return act
}

func (state *MobActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
