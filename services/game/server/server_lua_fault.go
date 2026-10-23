package server

import (
	"log"
	"time"

	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/fault"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/core/rpc"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game"
	lua "github.com/yuin/gopher-lua"
)

func registerFaultLuaFuncs(gs *GameServer, luaState *lua.LState) {
	injectors := map[string]*fault.Injector{
		"actor": gs.systems.Get[*c_actor.ActorRegistry]().Faults,
		"rpc":   gs.systems.Get[*rpc.InternalConnection]().Faults,
	}
	for name, in := range injectors {
		luax.RegisterFunc(luaState, name+"_delay", func(L *lua.LState) int {
			if L.Get(1).Type() == lua.LTString {
				key := L.CheckString(1)
				if L.GetTop() >= 2 {
					in.SetKeyDelay(key, time.Duration(L.CheckInt64(2))*time.Millisecond)
				}
				L.Push(lua.LNumber(in.KeyDelay(key).Milliseconds()))
				return 1
			}

			if L.GetTop() >= 1 {
				in.SetDelay(time.Duration(L.CheckInt64(1)) * time.Millisecond)
			}
			L.Push(lua.LNumber(in.Delay().Milliseconds()))
			return 1
		})

		luax.RegisterFunc(luaState, name+"_fault", func(L *lua.LState) int {
			key := L.CheckString(1)
			if L.GetTop() >= 2 {
				switch L.CheckString(2) {
				case "off":
					in.SetMode(key, fault.None)
				case "unreachable":
					in.SetMode(key, fault.Unreachable)
				case "lost":
					in.SetMode(key, fault.Lost)
				case "timeout":
					in.SetMode(key, fault.Timeout)
				default:
					L.ArgError(2, "off, unreachable, lost, timeout")
					return 0
				}
			}
			L.Push(lua.LString(in.Mode(key).String()))
			return 1
		})
	}

	luax.RegisterFunc(luaState, "settle_timer", func(L *lua.LState) int {
		if L.GetTop() >= 1 {
			gs.settlePaused.Store(!L.CheckBool(1))
		}
		L.Push(lua.LBool(!gs.settlePaused.Load()))
		return 1
	})

	luax.RegisterFunc(luaState, "mq_fault", func(L *lua.LState) int {
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LFalse)
			return 1
		}
		promise := gs.systems.Get[*game.InternalSystem]().SetPublishFaultAsync(cfg.ActorContext, L.CheckBool(1))
		return game.LuaYieldPromise(L, gs, promise, func(_ *internal.SetPublishFaultReply, err error) []lua.LValue {
			if err != nil {
				log.Printf("mq fault: %v", err)
				return []lua.LValue{lua.LFalse}
			}
			return []lua.LValue{lua.LTrue}
		})
	})
}
