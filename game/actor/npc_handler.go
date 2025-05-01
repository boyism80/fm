package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
)

func RegisterNpcHandlers(ctx protoactor.Context, m *NpcActor, h *handler.MessageHandler) {
	RegisterObjectHandlers(ctx, &m.Npc.Object, h)
}
