package handler

import (
	"github.com/boyism80/fm/model"
)

func RegisterNpcHandlers(m *model.Npc, handler *MessageHandler) {
	RegisterObjectHandlers(&m.Object, handler)
}
