package model

import (
	"github.com/boyism80/fm/handler"
)

func RegisterItemHandlers(m *ItemActor, h *handler.MessageHandler) {
	RegisterObjectHandlers(&m.Object, h)
}
