package entity

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (n *Npc) LuaTypeName() string {
	return "LuaNpc"
}

func (n *Npc) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			npc, ok := ud.Value.(*Npc)
			if !ok {
				L.ArgError(1, "Npc expected")
				return 0
			}
			if npc.Wz != nil && npc.Wz.BaseSpawn != nil {
				L.Push(lua.LNumber(npc.Wz.BaseSpawn.ID))
			} else {
				L.Push(lua.LNumber(0))
			}
			return 1
		},
		"show_effect": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			npc, ok := ud.Value.(*Npc)
			if !ok {
				L.ArgError(1, "Npc expected")
				return 0
			}
			action := L.CheckString(2)
			npc.ShowEffect(action)
			return 0
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			npc, ok := ud.Value.(*Npc)
			if !ok {
				L.ArgError(1, "Npc expected")
				return 0
			}
			if npc.Wz == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, npc.Wz))
			return 1
		},
	}
}

func (n *Npc) String() string {
	return n.LuaTypeName()
}

func (n *Npc) Type() lua.LValueType {
	return lua.LTUserData
}
