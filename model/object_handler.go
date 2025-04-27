package model

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/dao"
	"github.com/boyism80/fm/handler"
	"github.com/boyism80/fm/msg"
)

func RegisterObjectHandlers(m *dao.Object, h *handler.MessageHandler) {
	handler.RegisterHandler(m, h, onObjectMove)
}

func onObjectMove(obj *dao.Object, ctx actor.Context, m *msg.ObjectMove) {
	obj.Position.X += m.X
	obj.Position.Y += m.Y

	fmt.Println("move")
}
