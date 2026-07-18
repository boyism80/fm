package entity

import (
	lua "github.com/yuin/gopher-lua"
)

func (p *Portal) LuaTypeName() string {
	return "LuaPortal"
}

func (p *Portal) String() string {
	return p.LuaTypeName()
}

func (p *Portal) Type() lua.LValueType {
	return lua.LTUserData
}

func (p *Portal) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			portal, ok := ud.Value.(*Portal)
			if !ok || portal == nil || portal.Wz == nil {
				L.ArgError(1, "Portal expected")
				return 0
			}
			L.Push(lua.LString(portal.Wz.Name))
			return 1
		},
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			portal, ok := ud.Value.(*Portal)
			if !ok || portal == nil || portal.Wz == nil {
				L.ArgError(1, "Portal expected")
				return 0
			}
			L.Push(lua.LNumber(portal.Wz.ID))
			return 1
		},
		"target": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			portal, ok := ud.Value.(*Portal)
			if !ok || portal == nil || portal.Wz == nil {
				L.ArgError(1, "Portal expected")
				return 0
			}
			L.Push(lua.LString(portal.Wz.Target))
			return 1
		},
		"target_map": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			portal, ok := ud.Value.(*Portal)
			if !ok || portal == nil || portal.Wz == nil {
				L.ArgError(1, "Portal expected")
				return 0
			}
			L.Push(lua.LNumber(portal.Wz.TargetMapId))
			return 1
		},
		"script": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			portal, ok := ud.Value.(*Portal)
			if !ok || portal == nil {
				L.ArgError(1, "Portal expected")
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LString(portal.Script()))
				return 1
			}
			portal.SetScript(L.CheckString(2))
			return 0
		},
	}
}
