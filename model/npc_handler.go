package model

import (
	"github.com/boyism80/fm/handler"
)

func RegisterNpcHandlers(m *NpcActor, h *handler.MessageHandler) {
	RegisterObjectHandlers(&m.Object, h)
}
