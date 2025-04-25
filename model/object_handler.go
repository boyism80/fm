package model

import (
	"fmt"

	protoactor "github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterObjectHandlers(m *Object, h *handler.MessageHandler) {
	handler.RegisterHandler(m, h, onObjectMove)
}

func onObjectMove(obj *Object, ctx protoactor.Context, m *msg.ObjectMove) {
	obj.Position.X += m.X
	obj.Position.Y += m.Y

	fmt.Println("move")
}
