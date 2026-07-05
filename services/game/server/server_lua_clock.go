package server

import (
	"log"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	g_actor "github.com/boyism80/fm/services/game/actor"
	lua "github.com/yuin/gopher-lua"
)

func registerClockLuaFuncs(gs *GameServer, luaState *lua.LState) {
	luax.RegisterFunc(luaState, "datetime", func(L *lua.LState) int {
		now := clock.Now()
		year, month, day, hour, minute, second := clock.DateTimeTable(now)
		tbl := L.NewTable()
		tbl.RawSetString("year", lua.LNumber(year))
		tbl.RawSetString("month", lua.LNumber(month))
		tbl.RawSetString("day", lua.LNumber(day))
		tbl.RawSetString("hour", lua.LNumber(hour))
		tbl.RawSetString("minute", lua.LNumber(minute))
		tbl.RawSetString("second", lua.LNumber(second))
		L.Push(tbl)
		return 1
	})

	luax.RegisterFunc(luaState, "now", func(L *lua.LState) int {
		argc := L.GetTop()
		if argc == 0 {
			L.Push(lua.LNumber(clock.Now().Unix()))
			return 1
		}
		if gs == nil || gs.internalClient == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("internal client unavailable"))
			return 2
		}

		reset := false
		datetime := ""
		if argc >= 1 {
			switch L.Get(1).Type() {
			case lua.LTString:
				value := L.CheckString(1)
				if value == "reset" {
					reset = true
				} else {
					datetime = value
				}
			default:
				L.Push(lua.LBool(false))
				L.Push(lua.LString("datetime string or reset required"))
				return 2
			}
		}
		if !reset && datetime == "" {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("datetime is required"))
			return 2
		}
		if !reset {
			if _, err := clock.ParseDateTime(datetime); err != nil {
				L.Push(lua.LBool(false))
				L.Push(lua.LString(err.Error()))
				return 2
			}
		}

		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("actor context not found"))
			return 2
		}
		actorCtx := cfg.ActorContext
		pid := actorCtx.Self()
		if pid == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("actor PID not found"))
			return 2
		}
		mapInstance := gs.getMapByActorPID(pid)
		if mapInstance == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("map not found"))
			return 2
		}
		root := mapInstance.GetLuaRoot()
		if root == nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("lua state not found"))
			return 2
		}
		cfg.KeepAlive = true
		luax.SetConfiguration(L, cfg)

		succeeded := true
		errMsg := ""
		gs.SetServerDateTimeAsync(actorCtx, reset, datetime).OnError(func(err error) {
			succeeded = false
			if err != nil {
				errMsg = err.Error()
				log.Printf("now: %v", err)
			}
		}).Finally(func() {
			args := []lua.LValue{lua.LBool(succeeded)}
			if errMsg != "" {
				args = append(args, lua.LString(errMsg))
			} else {
				args = append(args, lua.LNil)
			}
			gs.GetRootContext().Send(pid, &g_actor.ResumeLua{Root: root, Thread: L, Args: args})
		})
		return L.Yield(lua.LNil, lua.LNil)
	})

	luax.RegisterFunc(luaState, "time_forward", func(L *lua.LState) int {
		return shiftServerDateTimeFromLua(gs, L, true)
	})
	luax.RegisterFunc(luaState, "time_backward", func(L *lua.LState) int {
		return shiftServerDateTimeFromLua(gs, L, false)
	})
}

func shiftServerDateTimeFromLua(gs *GameServer, L *lua.LState, forward bool) int {
	if gs == nil || gs.internalClient == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("internal client unavailable"))
		return 2
	}
	raw := L.CheckString(1)
	delta, err := clock.ParseTimespan(raw)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	if !forward {
		delta = -delta
	}
	target := clock.Now().Add(delta)
	return requestServerDateTimeFromLua(gs, L, false, clock.FormatDateTime(target))
}

func requestServerDateTimeFromLua(gs *GameServer, L *lua.LState, reset bool, datetime string) int {
	cfg, ok := luax.GetConfiguration(L)
	if !ok || cfg.ActorContext == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("actor context not found"))
		return 2
	}
	actorCtx := cfg.ActorContext
	pid := actorCtx.Self()
	if pid == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("actor PID not found"))
		return 2
	}
	mapInstance := gs.getMapByActorPID(pid)
	if mapInstance == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("map not found"))
		return 2
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("lua state not found"))
		return 2
	}
	cfg.KeepAlive = true
	luax.SetConfiguration(L, cfg)

	succeeded := true
	errMsg := ""
	gs.SetServerDateTimeAsync(actorCtx, reset, datetime).OnError(func(err error) {
		succeeded = false
		if err != nil {
			errMsg = err.Error()
		}
	}).Finally(func() {
		args := []lua.LValue{lua.LBool(succeeded)}
		if errMsg != "" {
			args = append(args, lua.LString(errMsg))
		} else {
			args = append(args, lua.LNil)
		}
		gs.GetRootContext().Send(pid, &g_actor.ResumeLua{Root: root, Thread: L, Args: args})
	})
	return L.Yield(lua.LNil, lua.LNil)
}
