package entity

import (
	"github.com/boyism80/fm/core/luax"
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
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			L.Push(lua.LString(machine.ID))
			return 1
		},
		"group": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			if machine.Group == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, machine.Group))
			return 1
		},
		"players": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
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
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			if machine.Party == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, machine.Party))
			return 1
		},
		"leader": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			if machine.Leader == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, machine.Leader))
			return 1
		},
		"scale_level": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			L.Push(lua.LNumber(machine.ScaleLevel))
			return 1
		},
		"enter_player": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			chUd := L.CheckUserData(2)
			ch, ok := chUd.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			machine.EnterPlayer(ch)
			return 0
		},
		"start": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			machine.Start()
			return 0
		},
		"unregister": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			chUd := L.CheckUserData(2)
			ch, ok := chUd.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			cfg, _ := luax.GetConfiguration(L)
			machine.LeavePlayer(cfg.ActorContext, ch, false)
			return 0
		},
		"finish": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			exitMapID := uint32(0)
			if L.GetTop() >= 2 {
				exitMapID = uint32(L.CheckInt(2))
			}
			cfg, _ := luax.GetConfiguration(L)
			machine.Finish(cfg.ActorContext, exitMapID)
			return 0
		},
		"start_timer": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
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
			return luaYieldPromise(L, machine.Group.GameWorld, promise)
		},
		"restart_timer": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
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
			return luaYieldPromise(L, machine.Group.GameWorld, promise)
		},
		"stop_timer": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			cfg, cfgOK := luax.GetConfiguration(L)
			if !cfgOK || cfg.ActorContext == nil {
				machine.StopTimer()
				return 0
			}
			promise := machine.StopTimerAsync(cfg.ActorContext)
			if promise == nil {
				return 0
			}
			return luaYieldPromise(L, machine.Group.GameWorld, promise)
		},
		"time_left": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			L.Push(lua.LNumber(machine.TimeLeft()))
			return 1
		},
		"set_property": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			key := L.CheckString(2)
			value := L.CheckString(3)
			machine.SetProperty(key, value)
			return 0
		},
		"get_property": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			key := L.CheckString(2)
			L.Push(lua.LString(machine.GetProperty(key)))
			return 1
		},
		"add_kill": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			machine, ok := ud.Value.(*StateMachine)
			if !ok || machine == nil {
				L.ArgError(1, "StateMachine expected")
				return 0
			}
			chUd := L.CheckUserData(2)
			ch, ok := chUd.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			n := 1
			if L.GetTop() >= 3 {
				n = L.CheckInt(3)
			}
			machine.AddKill(ch, n)
			return 0
		},
	}
}
