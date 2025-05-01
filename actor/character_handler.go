package actor

import (
	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
)

func RegisterCharacterHandlers(ctx protoactor.Context, m *entity.Character, h *handler.MessageHandler) {
	RegisterLifeHandlers(ctx, &m.Life, h)
}
