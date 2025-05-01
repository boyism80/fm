package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
)

func RegisterMobHandlers(ctx actor.Context, m *MobActor, h *handler.MessageHandler) {
	RegisterLifeHandlers(ctx, &m.Mob.Life, h)
}
