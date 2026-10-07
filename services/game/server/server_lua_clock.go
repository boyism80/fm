package server

import (
	"log"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/entity"
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

		return requestServerDateTimeFromLua(gs, L, reset, datetime)
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
	return entity.LuaYieldPromise(L, gs, gs.SetServerDateTimeAsync(cfg.ActorContext, reset, datetime), func(_ *internal.SetServerDateTimeReply, err error) []lua.LValue {
		if err != nil {
			log.Printf("server datetime: %v", err)
			return []lua.LValue{lua.LFalse, lua.LString(err.Error())}
		}
		return []lua.LValue{lua.LTrue, lua.LNil}
	})
}
