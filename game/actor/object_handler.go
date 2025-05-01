package actor

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/msg"
)

func RegisterObjectHandlers(ctx protoactor.Context, m *entity.Object, h *handler.MessageHandler) {
	handler.RegisterHandler(ctx, m, h, onObjectMove)
}

func onObjectMove(ctx protoactor.Context, obj *entity.Object, m *msg.ObjectMove) {
	obj.Position.X += m.X
	obj.Position.Y += m.Y

	fmt.Println("move")
}
