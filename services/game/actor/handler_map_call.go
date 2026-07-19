package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type MapCallHandler struct{}

func (MapCallHandler) New() *MapCallHandler {
	return &MapCallHandler{}
}

func (h *MapCallHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *MapCall) {
	if msg == nil || msg.Run == nil {
		return
	}
	msg.Values = msg.Run(ctx, a)
	if msg.ReplyTo == nil {
		return
	}
	ctx.Send(msg.ReplyTo, &MapCallAck{
		Root:   msg.Root,
		Thread: msg.Thread,
		Values: msg.Values,
	})
}
