package wz

import lua "github.com/yuin/gopher-lua"

func (*Pet) LuaTypeName() string { return luaWzPetTypeName }
func (*Pet) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			pushWzItemID(L, ud.Value)
			return 1
		},
	}
}

func (p *Pet) String() string       { return p.LuaTypeName() }
func (p *Pet) Type() lua.LValueType { return lua.LTUserData }
