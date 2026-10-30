package luax

import (
	"strings"

	"github.com/boyism80/fm/types"
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

// ParsePoint reads {x, y} or {x = x, y = y}.
func ParsePoint(table *lua.LTable) types.Point[int16] {
	var point types.Point[int16]
	if x := table.RawGetInt(1); x != lua.LNil {
		point.X = int16(lua.LVAsNumber(x))
	} else if x := table.RawGetString("x"); x != lua.LNil {
		point.X = int16(lua.LVAsNumber(x))
	}
	if y := table.RawGetInt(2); y != lua.LNil {
		point.Y = int16(lua.LVAsNumber(y))
	} else if y := table.RawGetString("y"); y != lua.LNil {
		point.Y = int16(lua.LVAsNumber(y))
	}
	return point
}
