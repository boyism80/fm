package luax

import (
	lua "github.com/yuin/gopher-lua"
)

func RegisterRequire(L *lua.LState) {
	L.SetGlobal("require", L.NewFunction(requireModule))
}

func requireModule(L *lua.LState) int {
	name := L.CheckString(1)
	path := name
	if len(path) < 4 || path[len(path)-4:] != ".lua" {
		path += ".lua"
	}

	mod, err := loadModule(L, path)
	if err != nil {
		L.RaiseError("require %s: %v", name, err)
		return 0
	}
	L.Push(mod)
	return 1
}
