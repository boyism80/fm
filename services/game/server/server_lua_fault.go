package server

import (
	"time"

	"github.com/boyism80/fm/core/fault"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func registerFaultLuaFuncs(gs *GameServer, luaState *lua.LState) {
	injectors := map[string]*fault.Injector{
		"actor": gs.actorRegistry.Faults,
		"rpc":   gs.rpcFaults,
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
}
