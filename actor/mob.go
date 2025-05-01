package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/context"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
)

type MobActor struct {
	entity.Mob
	handler *handler.MessageHandler
}

func NewMobActor(ctx protoactor.Context, serverCtx context.IServerContext) protoactor.Actor {
	act := &MobActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterMobHandlers(ctx, act, act.handler)
	return act
}

func (state *MobActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
