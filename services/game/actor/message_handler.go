package actor

import (
	"reflect"

	"github.com/asynkron/protoactor-go/actor"
)

type MessageHandler[M any] interface {
	Handle(ctx actor.Context, a *GameLogicActor, msg *M)
}

type MessageHandlerConstructor[H MessageHandler[M], M any] interface {
	New() H
}

type MessageHandlerFunc func(ctx actor.Context, a *GameLogicActor, msg any)

type MessageRegistry struct {
	byType map[reflect.Type]MessageHandlerFunc
}

func NewMessageRegistry() *MessageRegistry {
	return &MessageRegistry{byType: make(map[reflect.Type]MessageHandlerFunc)}
}

func Bind[C MessageHandlerConstructor[H, M], H MessageHandler[M], M any](r *MessageRegistry) {
	if r == nil {
		return
	}
	var constructor C
	handler := constructor.New()
	t := reflect.TypeOf((*M)(nil))
	r.byType[t] = func(ctx actor.Context, a *GameLogicActor, msg any) {
		handler.Handle(ctx, a, msg.(*M))
	}
}

func (r *MessageRegistry) Dispatch(ctx actor.Context, a *GameLogicActor, msg any) {
	if r == nil || msg == nil {
		return
	}
	h := r.byType[reflect.TypeOf(msg)]
	if h == nil {
		return
	}
	h(ctx, a, msg)
}

var mapMessageRegistry = NewMessageRegistry()

func init() {
	registerMapMessageHandlers(mapMessageRegistry)
}
