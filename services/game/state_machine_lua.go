package game

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	lua "github.com/yuin/gopher-lua"
)

func (sm *StateMachine) LuaTypeName() string {
	return "LuaStateMachine"
}

func (sm *StateMachine) String() string {
	return sm.LuaTypeName()
}

func (sm *StateMachine) Type() lua.LValueType {
	return lua.LTUserData
}

func (sm *StateMachine) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			L.Push(lua.LString(machine.id))
			return 1
		},
		"group": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			if machine.group == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, machine.group))
			return 1
		},
		"players": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			tbl := L.NewTable()
			for i, ch := range machine.Players() {
				if ch == nil {
					continue
				}
				tbl.RawSetInt(i+1, luax.NewLuable(L, ch))
			}
			L.Push(tbl)
			return 1
		},
		"party": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			if machine.party == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, machine.party))
			return 1
		},
		"leader": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			if machine.leader == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, machine.leader))
			return 1
		},
		"scale_level": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			L.Push(lua.LNumber(machine.scaleLevel))
			return 1
		},
		"disposed": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			L.Push(lua.LBool(machine.Disposed()))
			return 1
		},
		"enter_player": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			ch := luax.Check[Character](L, 2)
			args := []interface{}{}
			for i := 3; i <= L.GetTop(); i++ {
				args = append(args, L.Get(i))
			}
			machine.EnterPlayer(ch, args...)
			return 0
		},
		"start": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			machine.Start()
			return 0
		},
		"unregister": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			ch := luax.Check[Character](L, 2)
			cfg, _ := luax.GetConfiguration(L)
			if cfg.ActorContext != nil && cfg.ActorContext.Self().Equal(machine.actorPID) {
				machine.LeavePlayer(cfg.ActorContext, ch, false, StateMachineLeaveExit)
			} else {
				machine.RequestLeave(ch, false, StateMachineLeaveExit)
			}
			return 0
		},
		"map": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			m := machine.Map(uint32(L.CheckInt(2)))
			if m == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, m))
			return 1
		},
		"maps": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			tbl := L.NewTable()
			for i, m := range machine.MapList() {
				if m == nil {
					continue
				}
				tbl.RawSetInt(i+1, luax.NewLuable(L, m))
			}
			L.Push(tbl)
			return 1
		},
		"finish": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			exitMapID := uint32(L.OptInt(2, 0))
			exitPortal := uint8(L.OptInt(3, 0))
			cfg, _ := luax.GetConfiguration(L)
			if cfg.ActorContext != nil && cfg.ActorContext.Self().Equal(machine.actorPID) {
				machine.Finish(cfg.ActorContext, exitMapID, exitPortal)
			} else {
				machine.RequestFinish(exitMapID, exitPortal)
			}
			return 0
		},
		"start_timer": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			ms := int64(L.CheckNumber(2))
			cfg, cfgOK := luax.GetConfiguration(L)
			if !cfgOK || cfg.ActorContext == nil {
				machine.StartTimer(ms)
				return 0
			}
			promise := machine.StartTimerAsync(cfg.ActorContext, ms)
			if promise == nil {
				return 0
			}
			return LuaYieldPromise(L, machine.group.gameWorld, promise, nil)
		},
		"restart_timer": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			ms := int64(L.CheckNumber(2))
			cfg, cfgOK := luax.GetConfiguration(L)
			if !cfgOK || cfg.ActorContext == nil {
				machine.StartTimer(ms)
				return 0
			}
			promise := machine.StartTimerAsync(cfg.ActorContext, ms)
			if promise == nil {
				return 0
			}
			return LuaYieldPromise(L, machine.group.gameWorld, promise, nil)
		},
		"stop_timer": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			cfg, cfgOK := luax.GetConfiguration(L)
			if !cfgOK || cfg.ActorContext == nil {
				machine.StopTimer()
				return 0
			}
			promise := machine.StopTimerAsync(cfg.ActorContext)
			if promise == nil {
				return 0
			}
			return LuaYieldPromise(L, machine.group.gameWorld, promise, nil)
		},
		"time_left": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			L.Push(lua.LNumber(machine.TimeLeft()))
			return 1
		},
		"set_property": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			key := L.CheckString(2)
			value := L.CheckString(3)
			machine.SetProperty(key, value)
			return 0
		},
		"get_property": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			key := L.CheckString(2)
			L.Push(lua.LString(machine.GetProperty(key)))
			return 1
		},
		"add_kill": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			ch := luax.Check[Character](L, 2)
			n := 1
			if L.GetTop() >= 3 {
				n = L.CheckInt(3)
			}
			machine.AddKill(ch, n)
			return 0
		},
		"message": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			message := L.CheckString(2)
			msgType := constant.MsgLightBlueText
			if L.GetTop() >= 3 {
				msgType = constant.ServerMessageType(L.CheckInt(3))
			}
			for _, ch := range machine.Players() {
				if ch == nil {
					continue
				}
				ch.listener.OnMessage(ch, msgType, message)
			}
			return 0
		},
		"call_hook": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			hook := L.CheckString(2)
			if hook == "" {
				L.ArgError(2, "hook name expected")
				return 0
			}
			args := []interface{}{}
			for i := 3; i <= L.GetTop(); i++ {
				args = append(args, L.Get(i))
			}
			machine.CallHook(hook, args...)
			return 0
		},
		"after": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			var id string
			var ms int64
			var hook string
			if L.GetTop() >= 4 {
				id = L.CheckString(2)
				ms = int64(L.CheckNumber(3))
				hook = L.CheckString(4)
			} else {
				ms = int64(L.CheckNumber(2))
				hook = L.CheckString(3)
				id = hook
			}
			cfg, cfgOK := luax.GetConfiguration(L)
			if !cfgOK || cfg.ActorContext == nil {
				machine.Schedule(id, ms, hook)
				return 0
			}
			promise := machine.ScheduleAsync(cfg.ActorContext, id, ms, hook)
			if promise == nil {
				return 0
			}
			return LuaYieldPromise(L, machine.group.gameWorld, promise, nil)
		},
		"cron": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			id := L.CheckString(2)
			expr := L.CheckString(3)
			hook := L.CheckString(4)
			cfg, cfgOK := luax.GetConfiguration(L)
			if !cfgOK || cfg.ActorContext == nil {
				machine.Cron(id, expr, hook)
				return 0
			}
			promise := machine.CronAsync(cfg.ActorContext, id, expr, hook)
			if promise == nil {
				return 0
			}
			return LuaYieldPromise(L, machine.group.gameWorld, promise, nil)
		},
		"cancel": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			id := ""
			if L.GetTop() >= 2 {
				id = L.CheckString(2)
			}
			cfg, cfgOK := luax.GetConfiguration(L)
			if !cfgOK || cfg.ActorContext == nil {
				machine.CancelSchedule(id)
				return 0
			}
			promise := machine.CancelScheduleAsync(cfg.ActorContext, id)
			if promise == nil {
				return 0
			}
			return LuaYieldPromise(L, machine.group.gameWorld, promise, nil)
		},
		"warp_all": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			fromMapID := uint32(L.CheckInt(2))
			toMapID := uint32(L.CheckInt(3))
			spawnPoint := uint8(0)
			if L.GetTop() >= 4 {
				spawnPoint = uint8(L.CheckInt(4))
			}
			cfg, _ := luax.GetConfiguration(L)
			machine.WarpAll(cfg.ActorContext, fromMapID, toMapID, spawnPoint)
			return 0
		},
		"broadcast_ship": func(L *lua.LState) int {
			machine := luax.Check[StateMachine](L, 1)
			mapID := uint32(L.CheckInt(2))
			effect := uint16(L.CheckInt(3))
			machine.BroadcastShip(mapID, effect)
			return 0
		},
	}
}
