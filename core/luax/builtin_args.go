package luax

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// ArgError raises a Lua error and does not return, so a caller never receives nil.
func Check[T any, P interface {
	*T
	LuaTypeName() string
}](L *lua.LState, n int) P {
	v, ok := L.CheckUserData(n).Value.(P)
	if !ok || v == nil {
		L.ArgError(n, strings.TrimPrefix(P(nil).LuaTypeName(), "Lua")+" expected")
		return nil
	}
	return v
}

func ParseTable(L *lua.LState, argument lua.LValue, argIndex int) (map[lua.LValue]lua.LValue, bool) {
	table, ok := argument.(*lua.LTable)
	if !ok {
		L.ArgError(argIndex, "table expected")
		return nil, false
	}
	out := make(map[lua.LValue]lua.LValue, table.Len())
	table.ForEach(func(key lua.LValue, value lua.LValue) {
		out[key] = value
	})
	return out, true
}
