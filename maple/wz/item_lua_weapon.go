package wz

import lua "github.com/yuin/gopher-lua"

func (*Weapon) LuaTypeName() string { return luaWzWeaponTypeName }
func (*Weapon) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
		"weapon_type": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWeaponType(L, ud.Value)
			return 1
		},
	}
}

func (w *Weapon) String() string       { return w.LuaTypeName() }
func (w *Weapon) Type() lua.LValueType { return lua.LTUserData }
