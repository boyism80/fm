package wz

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

const luaWzNpcSpawnTypeName = "LuaWzNpcSpawn"

func (*NpcSpawn) LuaTypeName() string { return luaWzNpcSpawnTypeName }

func (*NpcSpawn) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			n, ok := ud.Value.(*NpcSpawn)
			if !ok || n == nil || n.BaseSpawn == nil {
				L.ArgError(1, "WzNpcSpawn expected")
				return 0
			}
			L.Push(lua.LNumber(n.BaseSpawn.ID))
			return 1
		},
		"position": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			n, ok := ud.Value.(*NpcSpawn)
			if !ok || n == nil || n.BaseSpawn == nil {
				L.ArgError(1, "WzNpcSpawn expected")
				return 0
			}
			posTbl := L.NewTable()
			posTbl.RawSetString("x", lua.LNumber(n.BaseSpawn.Position.X))
			posTbl.RawSetString("y", lua.LNumber(n.BaseSpawn.Position.Y))
			L.Push(posTbl)
			return 1
		},
		"collision_y": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			n, ok := ud.Value.(*NpcSpawn)
			if !ok || n == nil || n.BaseSpawn == nil {
				L.ArgError(1, "WzNpcSpawn expected")
				return 0
			}
			L.Push(lua.LNumber(n.BaseSpawn.CollisionY))
			return 1
		},
	}
}

func (n *NpcSpawn) String() string       { return n.LuaTypeName() }
func (n *NpcSpawn) Type() lua.LValueType { return lua.LTUserData }

var _ luax.Luable = (*NpcSpawn)(nil)
