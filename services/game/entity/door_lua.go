package entity

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (d *Door) LuaTypeName() string {
	return "LuaDoor"
}

func (d *Door) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"owner_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			door, ok := ud.Value.(*Door)
			if !ok {
				L.ArgError(1, "Door expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "owner_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(door.OwnerID))
			return 1
		},
		"skill_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			door, ok := ud.Value.(*Door)
			if !ok {
				L.ArgError(1, "Door expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "skill_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(door.SkillID))
			return 1
		},
		"remove": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			door, ok := ud.Value.(*Door)
			if !ok {
				L.ArgError(1, "Door expected")
				return 0
			}
			if door.Map != nil && door.OID != 0 {
				door.Map.RemoveDoor(door.OID, true)
			}
			return 0
		},
	}
}

var _ luax.Luable = (*Door)(nil)
