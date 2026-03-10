package entity

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (c *ItemCore) LuaTypeName() string { return "LuaItemCore" }

func (c *ItemCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			item, ok := ud.Value.(Item)
			if !ok {
				L.ArgError(1, "Item expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(item.GetModel().GetID()))
			return 1
		},
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
		"drop": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			item, ok := ud.Value.(Item)
			if !ok {
				L.ArgError(1, "Item expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "drop() is read-only")
				return 0
			}
			drop := item.GetDrop()
			if drop == nil || drop.Object == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, drop))
			return 1
		},
	}
}

func (c *ItemCore) String() string       { return c.LuaTypeName() }
func (c *ItemCore) Type() lua.LValueType { return lua.LTUserData }

func (*Equipment) LuaTypeName() string    { return "LuaEquipment" }
func (*Consume) LuaTypeName() string      { return "LuaConsume" }
func (*CashItem) LuaTypeName() string     { return "LuaCashItem" }
func (*GeneralItem) LuaTypeName() string  { return "LuaGeneralItem" }
func (*Installation) LuaTypeName() string { return "LuaInstallation" }
func (*Pet) LuaTypeName() string          { return "LuaPet" }

func (*Equipment) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*Consume) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*CashItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*GeneralItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*Installation) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*Pet) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}

func (e *Equipment) String() string    { return e.LuaTypeName() }
func (c *Consume) String() string      { return c.LuaTypeName() }
func (c *CashItem) String() string     { return c.LuaTypeName() }
func (g *GeneralItem) String() string  { return g.LuaTypeName() }
func (i *Installation) String() string { return i.LuaTypeName() }
func (p *Pet) String() string          { return p.LuaTypeName() }

func (e *Equipment) Type() lua.LValueType    { return lua.LTUserData }
func (c *Consume) Type() lua.LValueType      { return lua.LTUserData }
func (c *CashItem) Type() lua.LValueType     { return lua.LTUserData }
func (g *GeneralItem) Type() lua.LValueType  { return lua.LTUserData }
func (i *Installation) Type() lua.LValueType { return lua.LTUserData }
func (p *Pet) Type() lua.LValueType          { return lua.LTUserData }

func (d *Drop) LuaTypeName() string { return "LuaDrop" }
func (d *Drop) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (d *Drop) String() string       { return d.LuaTypeName() }
func (d *Drop) Type() lua.LValueType { return lua.LTUserData }

func (m *Meso) LuaTypeName() string { return "LuaMeso" }
func (m *Meso) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"count": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			meso, ok := ud.Value.(*Meso)
			if !ok {
				L.ArgError(1, "Meso expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "count() is read-only")
				return 0
			}
			L.Push(lua.LNumber(meso.Count))
			return 1
		},
	}
}
func (m *Meso) String() string       { return m.LuaTypeName() }
func (m *Meso) Type() lua.LValueType { return lua.LTUserData }

var (
	_ luax.Luable = (*ItemCore)(nil)
	_ luax.Luable = (*Drop)(nil)
	_ luax.Luable = (*Meso)(nil)
	_ luax.Luable = (*Equipment)(nil)
	_ luax.Luable = (*Consume)(nil)
	_ luax.Luable = (*CashItem)(nil)
	_ luax.Luable = (*GeneralItem)(nil)
	_ luax.Luable = (*Installation)(nil)
	_ luax.Luable = (*Pet)(nil)
)
