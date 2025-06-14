package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

type LuaRun struct {
	PID      *actor.PID
	FileName string
	FuncName string
	Params   []lua.LValue
}

type LuaResume struct {
	PID    *actor.PID
	Lua    *lua.LState
	Params []lua.LValue
}

type LuaYield struct {
	Lua *lua.LState
}
