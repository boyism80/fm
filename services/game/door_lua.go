package game

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
			door := luax.Check[Door](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "owner_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(door.ownerID))
			return 1
		},
		"skill_id": func(L *lua.LState) int {
			door := luax.Check[Door](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "skill_id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(door.skillID))
			return 1
		},
		"remove": func(L *lua.LState) int {
			door := luax.Check[Door](L, 1)
			if door.onMap != nil && door.oid != 0 {
				door.gameWorld.Systems().Get[*MapSystem]().DespawnDoor(door.onMap, door.key, true, true)
			}
			return 0
		},
	}
}

var _ luax.Luable = (*Door)(nil)
