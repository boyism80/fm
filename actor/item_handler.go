package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
)

func RegisterItemHandlers(ctx actor.Context, m *ItemActor, h *handler.MessageHandler) {
	RegisterObjectHandlers(ctx, &m.Item.Object, h)
}
