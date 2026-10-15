package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
	lua "github.com/yuin/gopher-lua"
)

type MapCallAsyncHandler struct{}

func (MapCallAsyncHandler) New() *MapCallAsyncHandler {
	return &MapCallAsyncHandler{}
}

func (h *MapCallAsyncHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.MapCallAsync) {
	if msg == nil || msg.Run == nil || msg.ReplyTo == nil {
		return
	}
	system := ctx.ActorSystem()
	if system == nil || system.Root == nil {
		return
	}
	promise := msg.Run(ctx, a.Maps())
	if promise == nil {
		system.Root.Send(msg.ReplyTo, &game.ResumeLua{
			Root:   msg.Root,
			Thread: msg.Thread,
			Args:   []lua.LValue{lua.LBool(false), lua.LNil, lua.LString("map call returned no promise")},
		})
		return
	}
	promise.Do(func(values []lua.LValue) error {
		system.Root.Send(msg.ReplyTo, &game.ResumeLua{
			Root:   msg.Root,
			Thread: msg.Thread,
			Args:   values,
		})
		return nil
	}).OnError(func(err error) {
		system.Root.Send(msg.ReplyTo, &game.ResumeLua{
			Root:   msg.Root,
			Thread: msg.Thread,
			Args:   []lua.LValue{lua.LBool(false), lua.LNil, lua.LString(err.Error())},
		})
	})
}

func (a *GameLogicActor) RunMapCallAsync(ctx actor.Context, msg *game.MapCallAsync) {
	(&MapCallAsyncHandler{}).Handle(ctx, a, msg)
}
