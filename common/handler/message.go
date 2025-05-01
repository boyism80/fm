package handler

import (
	"reflect"

	"github.com/asynkron/protoactor-go/actor"
)

type MessageHandler struct {
	handlers map[reflect.Type]func(ctx actor.Context, msg interface{})
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		handlers: make(map[reflect.Type]func(ctx actor.Context, msg interface{})),
	}
}

func (m *MessageHandler) Register(msgType reflect.Type, handler func(ctx actor.Context, msg interface{})) {
	m.handlers[msgType] = handler
}

func (m *MessageHandler) Handle(ctx actor.Context) {
	msg := ctx.Message()
	msgType := reflect.TypeOf(msg)

	if handler, ok := m.handlers[msgType]; ok {
		handler(ctx, msg)
	}
}

func RegisterHandler[T any, M any](ctx actor.Context, target *T, h *MessageHandler, fn func(actor.Context, *T, *M)) {
	var msg *M
	h.Register(reflect.TypeOf(msg), func(ctx actor.Context, m interface{}) {
		fn(ctx, target, m.(*M))
	})
}
