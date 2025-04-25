package actor

import (
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
)

func RegisterMobHandlers(m *model.Mob, h *handler.MessageHandler) {
	RegisterLifeHandlers(&m.Life, h)
}
