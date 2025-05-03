package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/game/entity"
)

type MobActor struct {
	entity.Mob
	handler *handler.MessageHandler
}

func NewMobActor(ctx protoactor.Context, serverCtx context.ServerContext) protoactor.Actor {
	act := &MobActor{
		handler: handler.NewMessageHandler(),
	}
	RegisterMobHandlers(ctx, act, act.handler)
	return act
}

func (state *MobActor) Receive(context protoactor.Context) {
	state.handler.Handle(context)
}
