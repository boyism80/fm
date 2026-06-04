package entity

import lua "github.com/yuin/gopher-lua"

func (a *Alliance) LuaTypeName() string {
	return "LuaAlliance"
}

func (a *Alliance) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			alliance, ok := ud.Value.(*Alliance)
			if !ok || alliance == nil {
				L.ArgError(1, "Alliance expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(alliance.AllianceID))
			return 1
		},
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			alliance, ok := ud.Value.(*Alliance)
			if !ok || alliance == nil {
				L.ArgError(1, "Alliance expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "name() is read-only")
				return 0
			}
			L.Push(lua.LString(alliance.Name))
			return 1
		},
		"capacity": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			alliance, ok := ud.Value.(*Alliance)
			if !ok || alliance == nil {
				L.ArgError(1, "Alliance expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "capacity() is read-only")
				return 0
			}
			L.Push(lua.LNumber(alliance.Capacity))
			return 1
		},
		"notice": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			alliance, ok := ud.Value.(*Alliance)
			if !ok || alliance == nil {
				L.ArgError(1, "Alliance expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "notice() is read-only")
				return 0
			}
			L.Push(lua.LString(alliance.Notice))
			return 1
		},
	}
}

func (a *Alliance) String() string {
	return a.LuaTypeName()
}

func (a *Alliance) Type() lua.LValueType {
	return lua.LTUserData
}
