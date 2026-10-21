package server

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game"
	lua "github.com/yuin/gopher-lua"
)

func registerFireWorksLuaFuncs(gs *GameServer, luaState *lua.LState) {
	luax.RegisterFunc(luaState, "fireworks_percentage", func(L *lua.LState) int {
		L.Push(lua.LNumber(gs.systems.Get[*game.FireWorksSystem]().Percentage(L.CheckString(1))))
		return 1
	})

	luax.RegisterFunc(luaState, "fireworks_give", func(L *lua.LState) int {
		L.Push(lua.LBool(gs.systems.Get[*game.FireWorksSystem]().Give(L.CheckString(1), L.CheckInt(2))))
		return 1
	})
}
