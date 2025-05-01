package actor

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
)

func RegisterLifeHandlers(ctx protoactor.Context, m *entity.Life, h *handler.MessageHandler) {
	RegisterObjectHandlers(ctx, &m.Object, h)

	handler.RegisterHandler(ctx, m, h, onLifeAddHp)
}

func onLifeAddHp(ctx protoactor.Context, life *entity.Life, m *msg.LifeAddHp) {
	life.HP += m.Hp

	fmt.Println("add hp")
}
