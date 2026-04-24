package entity

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (c *ItemCore) LuaTypeName() string { return "LuaItemCore" }

func (c *ItemCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	fieldPlacementGetter := func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		item, ok := ud.Value.(Item)
		if !ok {
			L.ArgError(1, "Item expected")
			return 0
		}
		if L.GetTop() != 1 {
			L.ArgError(2, "field_placement() is read-only")
			return 0
		}
		fp := item.GetFieldPlacement()
		if fp == nil || fp.ObjectCore == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(luax.NewLuable(L, fp))
		return 1
	}
	return map[string]lua.LGFunction{
		"count": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			item, ok := ud.Value.(Item)
			if !ok {
				L.ArgError(1, "Item expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "count() is read-only")
				return 0
			}
			L.Push(lua.LNumber(item.GetCount()))
			return 1
		},
		"field_placement": fieldPlacementGetter,
		"drop":            fieldPlacementGetter,
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			item, ok := ud.Value.(Item)
			if !ok {
				L.ArgError(1, "Item expected")
				return 0
			}
			model := item.GetModel()
			if luable, ok := model.(luax.Luable); ok && luable != nil {
				L.Push(luax.NewLuable(L, luable))
				return 1
			}
			L.Push(lua.LNil)
			return 1
		},
	}
}

func (c *ItemCore) String() string       { return c.LuaTypeName() }
func (c *ItemCore) Type() lua.LValueType { return lua.LTUserData }
