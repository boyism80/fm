package game

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (e *Expedition) LuaTypeName() string { return "LuaExpedition" }

func (e *Expedition) String() string { return e.LuaTypeName() }

func (e *Expedition) Type() lua.LValueType { return lua.LTUserData }

func (e *Expedition) LuaBuiltinFuncs() map[string]lua.LGFunction {
	result := func(L *lua.LState, err error) int {
		if err != nil {
			L.Push(lua.LFalse)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(lua.LTrue)
		return 1
	}
	return map[string]lua.LGFunction{
		"name": func(L *lua.LState) int {
			L.Push(lua.LString(luax.Check[Expedition](L, 1).name))
			return 1
		},
		"leader": func(L *lua.LState) int {
			L.Push(luax.NewLuable(L, luax.Check[Expedition](L, 1).leader))
			return 1
		},
		"members": func(L *lua.LState) int {
			members := luax.Check[Expedition](L, 1).Members()
			tbl := L.NewTable()
			for i := range members {
				tbl.Append(luax.NewLuable(L, &members[i]))
			}
			L.Push(tbl)
			return 1
		},
		"banned": func(L *lua.LState) int {
			banned := luax.Check[Expedition](L, 1).Banned()
			tbl := L.NewTable()
			for i := range banned {
				tbl.Append(luax.NewLuable(L, &banned[i]))
			}
			L.Push(tbl)
			return 1
		},
		"role": func(L *lua.LState) int {
			e := luax.Check[Expedition](L, 1)
			L.Push(lua.LNumber(e.Role(luax.Check[Character](L, 2))))
			return 1
		},
		"time_left": func(L *lua.LState) int {
			L.Push(lua.LNumber(luax.Check[Expedition](L, 1).TimeLeft()))
			return 1
		},
		"join": func(L *lua.LState) int {
			e := luax.Check[Expedition](L, 1)
			return result(L, e.Join(luax.Check[Character](L, 2)))
		},
		"leave": func(L *lua.LState) int {
			e := luax.Check[Expedition](L, 1)
			return result(L, e.Leave(luax.Check[Character](L, 2)))
		},
		"kick": func(L *lua.LState) int {
			e := luax.Check[Expedition](L, 1)
			return result(L, e.Kick(luax.Check[Character](L, 2), uint32(L.CheckInt(3))))
		},
		"allow": func(L *lua.LState) int {
			e := luax.Check[Expedition](L, 1)
			return result(L, e.Allow(luax.Check[Character](L, 2), uint32(L.CheckInt(3))))
		},
		"start": func(L *lua.LState) int {
			e := luax.Check[Expedition](L, 1)
			sm, skipped, err := e.Start(luax.Check[Character](L, 2))
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			tbl := L.NewTable()
			for _, ch := range skipped {
				tbl.Append(luax.NewLuable(L, ch))
			}
			L.Push(luax.NewLuable(L, sm))
			L.Push(tbl)
			return 2
		},
	}
}

func (member *ExpeditionMember) LuaTypeName() string { return "LuaExpeditionMember" }

func (member *ExpeditionMember) String() string { return member.LuaTypeName() }

func (member *ExpeditionMember) Type() lua.LValueType { return lua.LTUserData }

func (member *ExpeditionMember) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			L.Push(lua.LNumber(luax.Check[ExpeditionMember](L, 1).id))
			return 1
		},
		"name": func(L *lua.LState) int {
			L.Push(lua.LString(luax.Check[ExpeditionMember](L, 1).name))
			return 1
		},
		"class": func(L *lua.LState) int {
			L.Push(lua.LNumber(luax.Check[ExpeditionMember](L, 1).class))
			return 1
		},
	}
}

func RegisterExpeditionLua(gw GameWorld, L *lua.LState) {
	luax.RegisterLuaType[*Expedition](L)
	luax.RegisterLuaType[*ExpeditionMember](L)

	roles := L.NewTable()
	roles.RawSetString("None", lua.LNumber(ExpeditionRoleNone))
	roles.RawSetString("Leader", lua.LNumber(ExpeditionRoleLeader))
	roles.RawSetString("Member", lua.LNumber(ExpeditionRoleMember))
	roles.RawSetString("Banned", lua.LNumber(ExpeditionRoleBanned))
	L.SetGlobal("ExpeditionRole", roles)

	errs := L.NewTable()
	errs.RawSetString("Closed", lua.LString(ErrExpeditionClosed.Error()))
	errs.RawSetString("Exists", lua.LString(ErrExpeditionExists.Error()))
	errs.RawSetString("Battling", lua.LString(ErrExpeditionBattling.Error()))
	errs.RawSetString("Full", lua.LString(ErrExpeditionFull.Error()))
	errs.RawSetString("Joined", lua.LString(ErrExpeditionJoined.Error()))
	errs.RawSetString("Reserved", lua.LString(ErrExpeditionReserved.Error()))
	errs.RawSetString("Banned", lua.LString(ErrExpeditionBanned.Error()))
	errs.RawSetString("NotMember", lua.LString(ErrExpeditionNotMember.Error()))
	errs.RawSetString("NotLeader", lua.LString(ErrExpeditionNotLeader.Error()))
	errs.RawSetString("Away", lua.LString(ErrExpeditionAway.Error()))
	errs.RawSetString("TooFew", lua.LString(ErrExpeditionTooFew.Error()))
	L.SetGlobal("ExpeditionError", errs)

	reg := gw.Systems().Get[*ExpeditionSystem]()
	tbl := L.NewTable()
	tbl.RawSetString("register", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		leader := luax.Check[Character](L, 2)
		opts := L.CheckTable(3)
		spec := ExpeditionSpec{
			group:      lua.LVAsString(opts.RawGetString("group")),
			minMembers: int(lua.LVAsNumber(opts.RawGetString("min_members"))),
			maxMembers: int(lua.LVAsNumber(opts.RawGetString("max_members"))),
			maxBattles: int(lua.LVAsNumber(opts.RawGetString("max_battles"))),
			recruitMs:  int64(lua.LVAsNumber(opts.RawGetString("recruit_ms"))),
			notice:     lua.LVAsString(opts.RawGetString("notice")),
		}
		if spec.group == "" || spec.maxMembers <= 0 || spec.maxBattles <= 0 || spec.recruitMs <= 0 {
			L.ArgError(3, "expedition spec needs group, max_members, max_battles, recruit_ms")
			return 0
		}

		e, err := reg.Register(name, leader, spec)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(luax.NewLuable(L, e))
		return 1
	}))
	tbl.RawSetString("find", L.NewFunction(func(L *lua.LState) int {
		e := reg.Find(L.CheckString(1))
		if e == nil {
			L.Push(lua.LNil)
			return 1
		}
		L.Push(luax.NewLuable(L, e))
		return 1
	}))
	tbl.RawSetString("battles", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(reg.Battles(L.CheckString(1))))
		return 1
	}))
	tbl.RawSetString("reserve", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		ch := luax.Check[Character](L, 2)
		L.Push(lua.LBool(reg.Reserve(name, ch)))
		return 1
	}))
	tbl.RawSetString("reservations", L.NewFunction(func(L *lua.LState) int {
		reservations := reg.Reservations(L.CheckString(1))
		out := L.NewTable()
		for i := range reservations {
			out.Append(luax.NewLuable(L, &reservations[i]))
		}
		L.Push(out)
		return 1
	}))
	L.SetGlobal("expedition", tbl)
}
