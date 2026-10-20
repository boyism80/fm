package server

import (
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game"
	lua "github.com/yuin/gopher-lua"
)

func registerMedalLuaFuncs(gs *GameServer, luaState *lua.LState) {
	luax.RegisterFunc(luaState, "medal_ranking", func(L *lua.LState) int {
		kind := L.CheckString(1)
		entries, valid := gs.systems.Get[*game.MedalSystem]().Entries(kind)
		tbl := L.NewTable()
		for _, entry := range entries {
			row := L.NewTable()
			row.RawSetString("character_id", lua.LNumber(entry.CharacterID))
			row.RawSetString("name", lua.LString(entry.Name))
			row.RawSetString("score", lua.LNumber(entry.Score))
			tbl.Append(row)
		}
		L.Push(tbl)
		L.Push(lua.LBool(valid))
		return 2
	})

	luax.RegisterFunc(luaState, "refresh_medal_rankings", func(L *lua.LState) int {
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LFalse)
			return 1
		}
		promise := gs.systems.Get[*game.MedalSystem]().RefreshAsync(cfg.ActorContext)
		return game.LuaYieldPromise(L, gs, promise, func(_ bool, err error) []lua.LValue {
			if err != nil {
				log.Printf("refresh medal rankings: %v", err)
				return []lua.LValue{lua.LFalse}
			}
			return []lua.LValue{lua.LTrue}
		})
	})

	luax.RegisterFunc(luaState, "recheck_medals", func(L *lua.LState) int {
		ch, ok := L.CheckUserData(1).Value.(*game.Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			return 0
		}
		gs.systems.Get[*game.MedalSystem]().Recheck(cfg.ActorContext, ch.GetID())
		return 0
	})

	luax.RegisterFunc(luaState, "submit_medal_score", func(L *lua.LState) int {
		ch, ok := L.CheckUserData(1).Value.(*game.Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}
		kind := L.CheckString(2)
		score := uint64(L.CheckInt64(3))
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LNil)
			return 1
		}
		promise := gs.systems.Get[*game.MedalSystem]().SubmitAsync(cfg.ActorContext, ch, kind, score)
		return game.LuaYieldPromise(L, gs, promise, func(rank uint32, err error) []lua.LValue {
			if err != nil {
				log.Printf("submit medal score %s: %v", kind, err)
				return []lua.LValue{lua.LNil}
			}
			return []lua.LValue{lua.LNumber(rank)}
		})
	})
}
