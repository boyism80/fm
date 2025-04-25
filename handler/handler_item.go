package handler

import (
	"github.com/boyism80/fm/model"
)

func RegisterItemHandlers(m *model.Item, handler *MessageHandler) {
	RegisterObjectHandlers(&m.Object, handler)
}
