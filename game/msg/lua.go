package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

type LuaRun struct {
	PID      *actor.PID
	FileName string
	FuncName string
	Params   []any
}

type LuaResume struct {
	PID    *actor.PID
	Lua    *lua.LState
	Params []any
}

type LuaYield struct {
	Lua *lua.LState
}
