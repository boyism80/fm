package actor

import (
	"github.com/boyism80/fm/handler"
)

func RegisterMobHandlers(m *MobActor, h *handler.MessageHandler) {
	RegisterLifeHandlers(&m.Mob.Life, h)
}
