package model

import (
	"github.com/boyism80/fm/handler"
)

func RegisterMobHandlers(m *MobActor, h *handler.MessageHandler) {
	RegisterLifeHandlers(&m.Life, h)
}
