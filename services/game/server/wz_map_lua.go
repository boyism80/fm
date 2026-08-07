package server

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

func registerWzMapCreateInstance(gs *GameServer, L *lua.LState) {
	var zero *wz.Map
	meta := L.GetTypeMetatable(zero.LuaTypeName()).(*lua.LTable)
	L.SetFuncs(meta, map[string]lua.LGFunction{
		"create_instance": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			tmpl, ok := ud.Value.(*wz.Map)
			if !ok || tmpl == nil {
				L.ArgError(1, "WzMap expected")
				return 0
			}
			opts := entity.DefaultMapInitOpts()
			if L.GetTop() >= 2 && L.Get(2).Type() == lua.LTTable {
				opts = entity.ParseMapInitOptsLua(L.CheckTable(2))
			}
			ms := gs.GetMapSystem()
			if ms == nil {
				L.Push(lua.LNil)
				L.Push(lua.LString("map system not ready"))
				return 2
			}
			m, err := ms.CreateInstanceMap(tmpl.ID, opts)
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
