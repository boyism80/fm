package handler

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/model"
	"github.com/boyism80/fm/msg"
)

func HandleObjectMove(obj *model.Object, ctx protoactor.Context, m interface{}) {
	move := m.(*msg.ObjectMove)
	obj.Position.X += move.X
	obj.Position.Y += move.Y

	fmt.Println("move")
}

func RegisterObjectHandlers(m *model.Object, handler *MessageHandler) {
	RegisterHandler(m, &msg.ObjectMove{}, handler, HandleObjectMove)
}
