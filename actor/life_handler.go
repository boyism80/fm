package actor

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterLifeHandlers(m *entity.Life, h *handler.MessageHandler) {
	RegisterObjectHandlers(&m.Object, h)

	handler.RegisterHandler(m, h, onLifeAddHp)
}

func onLifeAddHp(life *entity.Life, ctx protoactor.Context, m *msg.LifeAddHp) {
	life.HP += m.Hp

	fmt.Println("add hp")
}
