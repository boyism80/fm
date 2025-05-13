package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/game/entity"
)

type NpcActor struct {
	entity.Npc
	handler *handler.MessageHandler
}

func NewNpcActor(ctx protoactor.Context, serverCtx context.ServerContext) protoactor.Actor {
	act := &NpcActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterObjectHandlers(ctx, &act.Npc.Object, act.handler)
	return act
}

func (state *NpcActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
