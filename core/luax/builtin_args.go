package luax

import lua "github.com/yuin/gopher-lua"

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
