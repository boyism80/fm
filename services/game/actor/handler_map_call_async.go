package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

type MapCallAsyncHandler struct{}

func (MapCallAsyncHandler) New() *MapCallAsyncHandler {
	return &MapCallAsyncHandler{}
}

func (h *MapCallAsyncHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *MapCallAsync) {
	if msg == nil || msg.Run == nil || msg.ReplyTo == nil {
		return
	}
	system := ctx.ActorSystem()
	if system == nil || system.Root == nil {
		return
	}
	promise := msg.Run(ctx, a)
	if promise == nil {
		system.Root.Send(msg.ReplyTo, &ResumeLua{
			Root:   msg.Root,
			Thread: msg.Thread,
			Args:   []lua.LValue{lua.LBool(false), lua.LNil, lua.LString("map call returned no promise")},
		})
		return
	}
	promise.Then(func(value interface{}) (interface{}, error) {
		values, _ := value.([]lua.LValue)
		system.Root.Send(msg.ReplyTo, &ResumeLua{
			Root:   msg.Root,
			Thread: msg.Thread,
			Args:   values,
		})
		return nil, nil
	}).OnError(func(err error) {
		system.Root.Send(msg.ReplyTo, &ResumeLua{
			Root:   msg.Root,
			Thread: msg.Thread,
			Args:   []lua.LValue{lua.LBool(false), lua.LNil, lua.LString(err.Error())},
		})
	})
}
