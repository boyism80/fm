package game

import (
	"log"

	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	lua "github.com/yuin/gopher-lua"
)

func luaYieldGuildRPC[T any](L *lua.LState, ch *Character, failCode int, result *int, promise *async.Promise[T]) int {
	if result == nil {
		L.Push(lua.LNumber(failCode))
		return 1
	}
	if promise == nil {
		L.Push(lua.LNumber(*result))
		return 1
	}
	if ch == nil || ch.gameWorld == nil {
		L.Push(lua.LNumber(failCode))
		return 1
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		L.Push(lua.LNumber(failCode))
		return 1
	}
	ch.dialog.Wait(L)
	thread := L
	promise.Finally(func() {
		if err := ch.dialog.Wake(thread, lua.LNumber(*result)); err != nil {
			log.Printf("guild lua resume: character=%d: %v", ch.GetID(), err)
		}
	})
	return L.Yield(lua.LNumber(0))
}

func luaGuildDisband(L *lua.LState, ch *Character) int {
	luax.Check[Guild](L, 1)
	if L.GetTop() != 2 {
		L.ArgError(3, "disband(character) takes exactly one argument")
		return 0
	}
	result := int(constant.GuildDisbandResultFailed)
	if ch == nil || ch.gameWorld == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	promise := ch.gameWorld.Systems().Get[*GuildSystem]().DisbandAsync(cfg.ActorContext, ch, &result)
	return luaYieldGuildRPC(L, ch, result, &result, promise)
}

func luaGuildIncCapacity(L *lua.LState, ch *Character, extendedCap bool) int {
	luax.Check[Guild](L, 1)
	if L.GetTop() != 3 {
		L.ArgError(4, "inc_capacity(character, extendedCap) takes exactly two arguments")
		return 0
	}
	result := int(constant.GuildIncreaseCapacityResultFailed)
	if ch == nil || ch.gameWorld == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	promise := ch.gameWorld.Systems().Get[*GuildSystem]().IncCapacityAsync(cfg.ActorContext, ch, extendedCap, &result)
	return luaYieldGuildRPC(L, ch, result, &result, promise)
}

func luaCreateAlliance(L *lua.LState, ch *Character) int {
	if ch == nil {
		return 0
	}
	if L.GetTop() != 2 {
		L.ArgError(3, "create_alliance(name) takes exactly one argument")
		return 0
	}
	allianceName := L.CheckString(2)
	result := int(constant.AllianceCreateResultFailed)
	if ch.gameWorld == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	promise := ch.gameWorld.Systems().Get[*GuildSystem]().CreateAllianceAsync(cfg.ActorContext, ch, allianceName, &result)
	return luaYieldGuildRPC(L, ch, result, &result, promise)
}

func luaIncAllianceCapacity(L *lua.LState, ch *Character) int {
	if ch == nil {
		return 0
	}
	if L.GetTop() != 1 {
		L.ArgError(1, "inc_alliance_capacity() takes no arguments")
		return 0
	}
	result := int(constant.AllianceIncreaseCapacityResultFailed)
	if ch.gameWorld == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	promise := ch.gameWorld.Systems().Get[*AllianceSystem]().IncCapacityAsync(cfg.ActorContext, ch, &result)
	return luaYieldGuildRPC(L, ch, result, &result, promise)
}

func luaDisbandAlliance(L *lua.LState, ch *Character) int {
	if ch == nil {
		return 0
	}
	if L.GetTop() != 1 {
		L.ArgError(1, "disband_alliance() takes no arguments")
		return 0
	}
	result := int(constant.AllianceDisbandResultFailed)
	if ch.gameWorld == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	promise := ch.gameWorld.Systems().Get[*GuildSystem]().DisbandAllianceAsync(cfg.ActorContext, ch, &result)
	return luaYieldGuildRPC(L, ch, result, &result, promise)
}
