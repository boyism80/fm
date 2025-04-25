package model

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterLifeHandlers(m *Life, h *handler.MessageHandler) {
	RegisterObjectHandlers(&m.Object, h)

	handler.RegisterHandler(m, h, onLifeAddHp)
}

func onLifeAddHp(life *Life, ctx actor.Context, m *msg.LifeAddHp) {
	life.HP += m.Hp

	fmt.Println("add hp")
}
