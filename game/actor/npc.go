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
	actor := &NpcActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterObjectHandlers(ctx, &actor.Npc.Object, actor.handler)
	return actor
}

func (state *NpcActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
