package handler

import (
	"reflect"

	protoactor "github.com/asynkron/protoactor-go/actor"
)

type MessageHandler struct {
	handlers map[reflect.Type]func(ctx protoactor.Context, msg interface{})
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		handlers: make(map[reflect.Type]func(ctx protoactor.Context, msg interface{})),
	}
}

func (m *MessageHandler) Register(msgType reflect.Type, handler func(ctx protoactor.Context, msg interface{})) {
	m.handlers[msgType] = handler
}

func (m *MessageHandler) Handle(ctx protoactor.Context) {
	msg := ctx.Message()
	msgType := reflect.TypeOf(msg)

	if handler, ok := m.handlers[msgType]; ok {
		handler(ctx, msg)
	}
}

func RegisterHandler[T any, M any](target *T, h *MessageHandler, fn func(*T, protoactor.Context, *M)) {
	var msg *M
	h.Register(reflect.TypeOf(msg), func(ctx protoactor.Context, m interface{}) {
		fn(target, ctx, m.(*M))
	})
}
