package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

type NpcActor struct {
	Npc     *model.Npc
	handler *handler.MessageHandler
}

func NewNpcActor(Npc *model.Npc) actor.Actor {
	act := &NpcActor{
		Npc:     Npc,
		handler: handler.NewMessageHandler(),
	}
	handler.RegisterNpcHandlers(act.Npc, act.handler)
	return act
}

func (state *NpcActor) Receive(context actor.Context) {
	state.handler.Handle(context)
}
