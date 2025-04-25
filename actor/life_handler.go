package actor

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
)

func RegisterLifeHandlers(m *model.Life, h *handler.MessageHandler) {
	RegisterObjectHandlers(&m.Object, h)

	handler.RegisterHandler(m, h, handleLifeAddHp)
}

func handleLifeAddHp(life *model.Life, ctx actor.Context, m *msg.LifeAddHp) {
	life.HP += m.Hp

	fmt.Println("add hp")
}
