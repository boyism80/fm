package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type MapCallAckHandler struct{}

func (MapCallAckHandler) New() *MapCallAckHandler {
	return &MapCallAckHandler{}
}

func (h *MapCallAckHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *MapCallAck) {
	if msg == nil || msg.Root == nil || msg.Thread == nil {
		return
	}
	handler := ResumeLuaHandler{}
	handler.Handle(ctx, a, &ResumeLua{
		Root:   msg.Root,
		Thread: msg.Thread,
		Args:   msg.Values,
	})
}
