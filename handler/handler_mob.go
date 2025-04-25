package handler

import (
	"github.com/boyism80/fm/model"
)

func RegisterMobHandlers(m *model.Mob, handler *MessageHandler) {
	RegisterLifeHandlers(&m.Life, handler)
}
