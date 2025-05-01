package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/context"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
)

type NpcActor struct {
	entity.Npc
	handler *handler.MessageHandler
}

func NewNpcActor(ctx protoactor.Context, serverCtx context.IServerContext) protoactor.Actor {
	act := &NpcActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterNpcHandlers(ctx, act, act.handler)
	return act
}

func (state *NpcActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
