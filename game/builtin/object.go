package builtin

import (
	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

type LuaObject struct {
	PID     *actor.PID
	Context *actor.RootContext
}

func (obj *LuaObject) LuaTypeName() string {
	return "LuaObject"
}

func (obj *LuaObject) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}

func (obj *LuaObject) String() string {
	return obj.LuaTypeName()
}

func (obj *LuaObject) Type() lua.LValueType {
	return lua.LTUserData
}
