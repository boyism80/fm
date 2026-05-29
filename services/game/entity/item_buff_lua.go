package entity

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (b *ItemBuff) LuaTypeName() string {
	return "LuaItemBuff"
}

func (b *ItemBuff) String() string {
	return b.LuaTypeName()
}

func (b *ItemBuff) Type() lua.LValueType {
	return lua.LTUserData
}

func (b *ItemBuff) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			ib, ok := ud.Value.(*ItemBuff)
			if !ok || ib == nil || ib.Wz == nil {
				L.ArgError(1, "ItemBuff expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "wz() is read-only")
				return 0
			}
			L.Push(luax.NewLuable(L, ib.Wz))
			return 1
		},
	}
}
