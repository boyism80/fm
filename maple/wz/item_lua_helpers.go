package wz

import (
	"github.com/boyism80/fm/maple/constant"
	lua "github.com/yuin/gopher-lua"
)

func pushWzItemID(L *lua.LState, v interface{}) bool {
	item, ok := v.(Item)
	if !ok || item == nil {
		L.Push(lua.LNil)
		return false
	}
	L.Push(lua.LNumber(item.GetID()))
	return true
}

func pushWeaponType(L *lua.LState, v interface{}) bool {
	w, ok := v.(*Weapon)
	if !ok || w == nil {
		L.Push(lua.LNil)
		return false
	}
	if w.WeaponType == constant.WeaponTypeNone {
		L.Push(lua.LNil)
		return false
	}
	L.Push(lua.LNumber(w.WeaponType))
	return true
}

func pushConsumeType(L *lua.LState, v interface{}) bool {
	c, ok := v.(*Consume)
	if !ok || c == nil {
		L.Push(lua.LNil)
		return false
	}
	if c.ConsumeType == constant.ConsumeTypeNone {
		L.Push(lua.LNil)
		return false
	}
	L.Push(lua.LNumber(c.ConsumeType))
	return true
}
