package actor

import (
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

func RegisterNpcHandlers(m *model.Npc, h *handler.MessageHandler) {
	RegisterObjectHandlers(&m.Object, h)
}
