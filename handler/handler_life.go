package handler

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
)

func RegisterLifeHandlers(m *model.Life, handler *MessageHandler) {
	RegisterObjectHandlers(&m.Object, handler)

	RegisterHandler(m, &msg.LifeAddHp{}, handler, HandleLifeAddHp)
}

func HandleLifeAddHp(life *model.Life, ctx actor.Context, m interface{}) {
	addHp := m.(*msg.LifeAddHp)
	life.HP += addHp.Hp

	fmt.Println("add hp")
}
