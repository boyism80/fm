package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

func (a *MapActor) RespawnCall(_ actor.Context, includeNegativeMobTime bool) []lua.LValue {
	if a.Map == nil {
		return []lua.LValue{lua.LNumber(0)}
	}
	return []lua.LValue{lua.LNumber(a.Map.Respawn(includeNegativeMobTime))}
}
