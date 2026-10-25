package server

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/services/game"
	lua "github.com/yuin/gopher-lua"
)

func registerWzMapCreateInstance(gs *GameServer, L *lua.LState) {
	var zero *wz.Map
	meta := L.GetTypeMetatable(zero.LuaTypeName()).(*lua.LTable)
	L.SetFuncs(meta, map[string]lua.LGFunction{
		"create_instance": func(L *lua.LState) int {
			tmpl := luax.Check[wz.Map](L, 1)
			opts := game.DefaultMapInitOpts()
			var slot lua.LValue = lua.LNil
			if L.GetTop() >= 2 && L.Get(2).Type() == lua.LTTable {
				tbl := L.CheckTable(2)
				opts = game.ParseMapInitOptsLua(tbl)
				slot = tbl.RawGetString("slot")
			}
			ms := gs.systems.Get[*game.MapSystem]()
			if ms == nil {
				L.Push(lua.LNil)
				L.Push(lua.LString("map system not ready"))
				return 2
			}
			var m *game.Map
			var err error
			if n, ok := slot.(lua.LNumber); ok {
				m, err = ms.SlotInstance(tmpl.ID, uint32(n))
			} else {
				m, err = ms.CreateInstanceMap(tmpl.ID, opts)
			}
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			L.Push(luax.NewLuable(L, m))
			return 1
		},
	})
}
