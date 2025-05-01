package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/context"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
)

type ItemActor struct {
	entity.Item
	ServerCtx context.IServerContext
	handler   *handler.MessageHandler
}

func NewItemActor(ctx protoactor.Context, serverCtx context.IServerContext) protoactor.Actor {
	act := &ItemActor{
		ServerCtx: serverCtx,
		handler:   handler.NewMessageHandler(),
	}
	RegisterItemHandlers(ctx, act, act.handler)
	return act
}

func (state *ItemActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
