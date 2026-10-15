package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type MapCallHandler struct{}

func (MapCallHandler) New() *MapCallHandler {
	return &MapCallHandler{}
}

func (h *MapCallHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.MapCall) {
	if msg == nil || msg.Run == nil {
		return
	}
	msg.Values = msg.Run(ctx, a.Maps())
	if msg.ReplyTo == nil {
		return
	}
	ctx.Send(msg.ReplyTo, &game.ResumeLua{
		Root:   msg.Root,
		Thread: msg.Thread,
		Args:   msg.Values,
	})
}

func (a *GameLogicActor) RunMapCall(ctx actor.Context, msg *game.MapCall) {
	(&MapCallHandler{}).Handle(ctx, a, msg)
}
