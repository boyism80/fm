package game

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	lua "github.com/yuin/gopher-lua"
)

// start returns the request, or nil and the code to answer with when it never went out.
func luaYieldGuildRPC(L *lua.LState, ch *Character, failCode int, start func(ctx actor.Context) (*async.Promise[int], int)) int {
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		L.Push(lua.LNumber(failCode))
		return 1
	}
	promise, result := start(cfg.ActorContext)
	if promise == nil {
		L.Push(lua.LNumber(result))
		return 1
	}
	ch.dialog.Wait(L)
	thread := L
	promise.Do(func(code int) error {
		result = code
		return nil
	}).Finally(func() {
		if err := ch.dialog.Wake(thread, lua.LNumber(result)); err != nil {
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
	return luaYieldGuildRPC(L, ch, int(constant.GuildDisbandResultFailed), func(ctx actor.Context) (*async.Promise[int], int) {
		return ch.gameWorld.Systems().Get[*GuildSystem]().DisbandAsync(ctx, ch)
	})
}

func luaGuildIncCapacity(L *lua.LState, ch *Character, extendedCap bool) int {
	luax.Check[Guild](L, 1)
	if L.GetTop() != 3 {
		L.ArgError(4, "inc_capacity(character, extendedCap) takes exactly two arguments")
		return 0
	}
	return luaYieldGuildRPC(L, ch, int(constant.GuildIncreaseCapacityResultFailed), func(ctx actor.Context) (*async.Promise[int], int) {
		return ch.gameWorld.Systems().Get[*GuildSystem]().IncCapacityAsync(ctx, ch, extendedCap)
	})
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
	return luaYieldGuildRPC(L, ch, int(constant.AllianceCreateResultFailed), func(ctx actor.Context) (*async.Promise[int], int) {
		return ch.gameWorld.Systems().Get[*GuildSystem]().CreateAllianceAsync(ctx, ch, allianceName)
	})
}

func luaIncAllianceCapacity(L *lua.LState, ch *Character) int {
	if ch == nil {
		return 0
	}
	if L.GetTop() != 1 {
		L.ArgError(1, "inc_alliance_capacity() takes no arguments")
		return 0
	}
	return luaYieldGuildRPC(L, ch, int(constant.AllianceIncreaseCapacityResultFailed), func(ctx actor.Context) (*async.Promise[int], int) {
		return ch.gameWorld.Systems().Get[*AllianceSystem]().IncCapacityAsync(ctx, ch)
	})
}

func luaDisbandAlliance(L *lua.LState, ch *Character) int {
	if ch == nil {
		return 0
	}
	if L.GetTop() != 1 {
		L.ArgError(1, "disband_alliance() takes no arguments")
		return 0
	}
	return luaYieldGuildRPC(L, ch, int(constant.AllianceDisbandResultFailed), func(ctx actor.Context) (*async.Promise[int], int) {
		return ch.gameWorld.Systems().Get[*GuildSystem]().DisbandAllianceAsync(ctx, ch)
	})
}
