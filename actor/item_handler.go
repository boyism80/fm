package actor

import (
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

func RegisterItemHandlers(m *model.Item, h *handler.MessageHandler) {
	RegisterObjectHandlers(&m.Object, h)
}
