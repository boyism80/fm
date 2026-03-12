package wz

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/constant"
	lua "github.com/yuin/gopher-lua"
)

type ItemWzCore struct {
	Model Item
}

type ItemWzEquipment struct {
	*ItemWzCore
}
type ItemWzWeapon struct {
	*ItemWzEquipment
}
type ItemWzArmor struct {
	*ItemWzEquipment
}
type ItemWzConsume struct {
	*ItemWzCore
}
type ItemWzPet struct {
	*ItemWzCore
}
type ItemWzGeneralItem struct {
	*ItemWzCore
}
type ItemWzCashItem struct {
	*ItemWzCore
}
type ItemWzInstallation struct {
	*ItemWzCore
}

func (*ItemWzCore) LuaTypeName() string { return "LuaItemWzCore" }

func (*ItemWzCore) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			model := getItemWzModel(L, 1)
			if model == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(model.GetID()))
			return 1
		},
	}
}

func (w *ItemWzCore) String() string       { return w.LuaTypeName() }
func (w *ItemWzCore) Type() lua.LValueType { return lua.LTUserData }

func getItemWzModel(L *lua.LState, idx int) Item {
	ud := L.CheckUserData(idx)
	if ud.Value == nil {
		return nil
	}
	switch w := ud.Value.(type) {
	case *ItemWzCore:
		return w.Model
	case *ItemWzEquipment:
		return w.ItemWzCore.Model
	case *ItemWzWeapon:
		return w.ItemWzEquipment.ItemWzCore.Model
	case *ItemWzArmor:
		return w.ItemWzEquipment.ItemWzCore.Model
	case *ItemWzConsume:
		return w.ItemWzCore.Model
	case *ItemWzPet:
		return w.ItemWzCore.Model
	case *ItemWzGeneralItem:
		return w.ItemWzCore.Model
	case *ItemWzCashItem:
		return w.ItemWzCore.Model
	case *ItemWzInstallation:
		return w.ItemWzCore.Model
	default:
		return nil
	}
}

func PushItemWz(L *lua.LState, model Item) {
	if model == nil {
		L.Push(lua.LNil)
		return
	}
	core := &ItemWzCore{Model: model}
	switch model.(type) {
	case *Weapon:
		L.Push(luax.NewLuable(L, &ItemWzWeapon{ItemWzEquipment: &ItemWzEquipment{ItemWzCore: core}}))
	case *Armor:
		L.Push(luax.NewLuable(L, &ItemWzArmor{ItemWzEquipment: &ItemWzEquipment{ItemWzCore: core}}))
	case *Consume:
		L.Push(luax.NewLuable(L, &ItemWzConsume{ItemWzCore: core}))
	case *Pet:
		L.Push(luax.NewLuable(L, &ItemWzPet{ItemWzCore: core}))
	case *GeneralItem:
		L.Push(luax.NewLuable(L, &ItemWzGeneralItem{ItemWzCore: core}))
	case *CashItem:
		L.Push(luax.NewLuable(L, &ItemWzCashItem{ItemWzCore: core}))
	case *Installation:
		L.Push(luax.NewLuable(L, &ItemWzInstallation{ItemWzCore: core}))
	default:
		L.Push(luax.NewLuable(L, core))
	}
}

func (*ItemWzEquipment) LuaTypeName() string    { return "LuaItemWzEquipment" }
func (*ItemWzWeapon) LuaTypeName() string       { return "LuaItemWzWeapon" }
func (*ItemWzArmor) LuaTypeName() string        { return "LuaItemWzArmor" }
func (*ItemWzConsume) LuaTypeName() string      { return "LuaItemWzConsume" }
func (*ItemWzPet) LuaTypeName() string          { return "LuaItemWzPet" }
func (*ItemWzGeneralItem) LuaTypeName() string  { return "LuaItemWzGeneralItem" }
func (*ItemWzCashItem) LuaTypeName() string     { return "LuaItemWzCashItem" }
func (*ItemWzInstallation) LuaTypeName() string { return "LuaItemWzInstallation" }

func (*ItemWzEquipment) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*ItemWzWeapon) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"weapon_type": func(L *lua.LState) int {
			model := getItemWzModel(L, 1)
			if model == nil {
				return 0
			}
			wt := constant.GetWeaponType(model.GetID())
			if wt == constant.WeaponTypeNone {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(wt))
			return 1
		},
	}
}
func (*ItemWzArmor) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*ItemWzConsume) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"consume_type": func(L *lua.LState) int {
			model := getItemWzModel(L, 1)
			if model == nil {
				return 0
			}
			ct := constant.GetConsumeType(model.GetID())
			if ct == constant.ConsumeTypeNone {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(ct))
			return 1
		},
	}
}
func (*ItemWzPet) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*ItemWzGeneralItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*ItemWzCashItem) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}
func (*ItemWzInstallation) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{}
}

func (w *ItemWzEquipment) String() string    { return w.LuaTypeName() }
func (w *ItemWzWeapon) String() string       { return w.LuaTypeName() }
func (w *ItemWzArmor) String() string        { return w.LuaTypeName() }
func (w *ItemWzConsume) String() string      { return w.LuaTypeName() }
func (w *ItemWzPet) String() string          { return w.LuaTypeName() }
func (w *ItemWzGeneralItem) String() string  { return w.LuaTypeName() }
func (w *ItemWzCashItem) String() string     { return w.LuaTypeName() }
func (w *ItemWzInstallation) String() string { return w.LuaTypeName() }

func (w *ItemWzEquipment) Type() lua.LValueType    { return lua.LTUserData }
func (w *ItemWzWeapon) Type() lua.LValueType       { return lua.LTUserData }
func (w *ItemWzArmor) Type() lua.LValueType        { return lua.LTUserData }
func (w *ItemWzConsume) Type() lua.LValueType      { return lua.LTUserData }
func (w *ItemWzPet) Type() lua.LValueType          { return lua.LTUserData }
func (w *ItemWzGeneralItem) Type() lua.LValueType  { return lua.LTUserData }
func (w *ItemWzCashItem) Type() lua.LValueType     { return lua.LTUserData }
func (w *ItemWzInstallation) Type() lua.LValueType { return lua.LTUserData }

var (
	_ luax.Luable = (*ItemWzCore)(nil)
	_ luax.Luable = (*ItemWzEquipment)(nil)
	_ luax.Luable = (*ItemWzWeapon)(nil)
	_ luax.Luable = (*ItemWzArmor)(nil)
	_ luax.Luable = (*ItemWzConsume)(nil)
	_ luax.Luable = (*ItemWzPet)(nil)
	_ luax.Luable = (*ItemWzGeneralItem)(nil)
	_ luax.Luable = (*ItemWzCashItem)(nil)
	_ luax.Luable = (*ItemWzInstallation)(nil)
)
