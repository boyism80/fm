package luax

import (
	lua "github.com/yuin/gopher-lua"
)

func RegisterRequire(L *lua.LState) {
	packageTable := L.GetGlobal("package")
	if packageTable == lua.LNil {
		packageTable = L.NewTable()
		L.SetGlobal("package", packageTable)
	}
	if L.GetField(packageTable, "loaded") == lua.LNil {
		L.SetField(packageTable, "loaded", L.NewTable())
	}

	L.SetGlobal("require", L.NewFunction(requireModule))
}

func clearRequireCache(L *lua.LState) {
	packageTable := L.GetGlobal("package")
	if packageTable == lua.LNil {
		return
	}
	pt, ok := packageTable.(*lua.LTable)
	if !ok {
		return
	}
	L.SetField(pt, "loaded", L.NewTable())
}

func requireModule(L *lua.LState) int {
	name := L.CheckString(1)
	path := name
	if len(path) < 4 || path[len(path)-4:] != ".lua" {
		path += ".lua"
	}

	packageTable := L.GetGlobal("package").(*lua.LTable)
	loaded := L.GetField(packageTable, "loaded").(*lua.LTable)

	compileMu.Lock()
	reload := alwaysReload
	compileMu.Unlock()

	if reload {
		L.SetField(loaded, name, lua.LNil)
	} else if cached := L.GetField(loaded, name); cached != lua.LNil {
		L.Push(cached)
		return 1
	}

	fn, err := L.LoadFile(path)
	if err != nil {
		L.RaiseError("require %s: %v", name, err)
		return 0
	}

	L.Push(fn)
	if err := L.PCall(0, 1, nil); err != nil {
		L.RaiseError("require %s: %v", name, err)
		return 0
	}

	mod := L.Get(-1)
	L.Pop(1)
	if mod == lua.LNil {
		mod = lua.LBool(true)
	}

	L.SetField(loaded, name, mod)
	L.Push(mod)
	return 1
}
