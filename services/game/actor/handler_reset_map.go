package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

func (a *MapActor) ResetCall(_ actor.Context) []lua.LValue {
	if a.Map == nil {
		return []lua.LValue{lua.LBool(false)}
	}
	a.Map.Reset()
	return []lua.LValue{lua.LBool(true)}
}
