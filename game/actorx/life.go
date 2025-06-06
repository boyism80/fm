package actorx

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
)

func RegisterLifeHandlers(ctx actor.Context, m *entity.Life, h *handler.MessageHandler) {
	RegisterObjectHandlers(ctx, &m.Object, h)

	handler.RegisterHandler(ctx, m, h, onLifeAddHp)
}

func onLifeAddHp(ctx actor.Context, life *entity.Life, m *msg.LifeAddHp) {
	life.Hp += uint16(m.Hp)

	fmt.Println("add hp")
}
