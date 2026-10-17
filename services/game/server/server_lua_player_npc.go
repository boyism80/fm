package server

import (
	"errors"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game"
	lua "github.com/yuin/gopher-lua"
)

func registerPlayerNpcLuaFuncs(gs *GameServer, luaState *lua.LState) {
	luax.RegisterFunc(luaState, "create_player_npc", func(L *lua.LState) int {
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString("no actor context"))
			return 2
		}
		at, ok := L.CheckUserData(1).Value.(*game.Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}
		npcID := uint32(L.CheckInt(2))
		source := at
		if L.GetTop() >= 3 {
			if source, ok = L.CheckUserData(3).Value.(*game.Character); !ok {
				L.ArgError(3, "Character expected")
				return 0
			}
		}

		npc, err := game.NewPlayerNpc(source, npcID, at)
		if err != nil {
			L.Push(lua.LFalse)
			switch {
			case errors.Is(err, game.ErrPlayerNpcUnknownTemplate):
				L.Push(lua.LString("unknown"))
			case errors.Is(err, game.ErrPlayerNpcTaken):
				L.Push(lua.LString("taken"))
			default:
				L.Push(lua.LString("placed"))
			}
			return 2
		}
		promise := gs.systems.Get[*game.PlayerNpcSystem]().CreateAsync(cfg.ActorContext, npc)
		return game.LuaYieldPromise(L, gs, promise, func(created bool, err error) []lua.LValue {
			if err != nil {
				log.Printf("create player npc %d: %v", npcID, err)
				return []lua.LValue{lua.LFalse, lua.LString("error")}
			}
			if !created {
				return []lua.LValue{lua.LFalse, lua.LString("taken")}
			}
			return []lua.LValue{lua.LTrue}
		})
	})

	luax.RegisterFunc(luaState, "remove_player_npc", func(L *lua.LState) int {
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LFalse)
			return 1
		}
		npcID := uint32(L.CheckInt(1))
		promise := gs.systems.Get[*game.PlayerNpcSystem]().RemoveAsync(cfg.ActorContext, npcID)
		return game.LuaYieldPromise(L, gs, promise, func(removed bool, err error) []lua.LValue {
			if err != nil {
				log.Printf("remove player npc %d: %v", npcID, err)
				return []lua.LValue{lua.LFalse}
			}
			return []lua.LValue{lua.LBool(removed)}
		})
	})
}
