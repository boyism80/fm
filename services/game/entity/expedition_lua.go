package entity

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (e *Expedition) LuaTypeName() string { return "LuaExpedition" }

func (e *Expedition) String() string { return e.LuaTypeName() }

func (e *Expedition) Type() lua.LValueType { return lua.LTUserData }

func (e *Expedition) LuaBuiltinFuncs() map[string]lua.LGFunction {
	check := func(L *lua.LState) *Expedition {
		ud := L.CheckUserData(1)
		e, ok := ud.Value.(*Expedition)
		if ok == false || e == nil {
			L.ArgError(1, "Expedition expected")
			return nil
		}
		return e
	}
	character := func(L *lua.LState, n int) *Character {
		ud := L.CheckUserData(n)
		ch, ok := ud.Value.(*Character)
		if ok == false || ch == nil {
			L.ArgError(n, "Character expected")
			return nil
		}
		return ch
	}
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
			L.Push(lua.LString(check(L).Name))
			return 1
		},
		"leader": func(L *lua.LState) int {
			L.Push(luax.NewLuable(L, check(L).Leader))
			return 1
		},
		"members": func(L *lua.LState) int {
			members := check(L).Members()
			tbl := L.NewTable()
			for i := range members {
				tbl.Append(luax.NewLuable(L, &members[i]))
			}
			L.Push(tbl)
			return 1
		},
		"banned": func(L *lua.LState) int {
			banned := check(L).Banned()
			tbl := L.NewTable()
			for i := range banned {
				tbl.Append(luax.NewLuable(L, &banned[i]))
			}
			L.Push(tbl)
			return 1
		},
		"role": func(L *lua.LState) int {
			e := check(L)
			L.Push(lua.LNumber(e.Role(character(L, 2))))
			return 1
		},
		"time_left": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).TimeLeft()))
			return 1
		},
		"join": func(L *lua.LState) int {
			e := check(L)
			return result(L, e.Join(character(L, 2)))
		},
		"leave": func(L *lua.LState) int {
			e := check(L)
			return result(L, e.Leave(character(L, 2)))
		},
		"kick": func(L *lua.LState) int {
			e := check(L)
			return result(L, e.Kick(character(L, 2), uint32(L.CheckInt(3))))
		},
		"allow": func(L *lua.LState) int {
			e := check(L)
			return result(L, e.Allow(character(L, 2), uint32(L.CheckInt(3))))
		},
		"start": func(L *lua.LState) int {
			e := check(L)
			sm, skipped, err := e.Start(character(L, 2))
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
	check := func(L *lua.LState) *ExpeditionMember {
		ud := L.CheckUserData(1)
		member, ok := ud.Value.(*ExpeditionMember)
		if ok == false || member == nil {
			L.ArgError(1, "ExpeditionMember expected")
			return nil
		}
		return member
	}
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).ID))
			return 1
		},
		"name": func(L *lua.LState) int {
			L.Push(lua.LString(check(L).Name))
			return 1
		},
		"class": func(L *lua.LState) int {
			L.Push(lua.LNumber(check(L).Class))
			return 1
		},
	}
}

func RegisterExpeditionLua(L *lua.LState, gw GameWorld) {
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

	reg := gw.GetExpeditionRegistry()
	tbl := L.NewTable()
	tbl.RawSetString("register", L.NewFunction(func(L *lua.LState) int {
		name := L.CheckString(1)
		ud := L.CheckUserData(2)
		leader, ok := ud.Value.(*Character)
		if ok == false || leader == nil {
			L.ArgError(2, "Character expected")
			return 0
		}
		opts := L.CheckTable(3)
		spec := ExpeditionSpec{
			Group:      lua.LVAsString(opts.RawGetString("group")),
			MinMembers: int(lua.LVAsNumber(opts.RawGetString("min_members"))),
			MaxMembers: int(lua.LVAsNumber(opts.RawGetString("max_members"))),
			MaxBattles: int(lua.LVAsNumber(opts.RawGetString("max_battles"))),
			RecruitMs:  int64(lua.LVAsNumber(opts.RawGetString("recruit_ms"))),
			Notice:     lua.LVAsString(opts.RawGetString("notice")),
		}
		if spec.Group == "" || spec.MaxMembers <= 0 || spec.MaxBattles <= 0 || spec.RecruitMs <= 0 {
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
		ud := L.CheckUserData(2)
		ch, ok := ud.Value.(*Character)
		if ok == false || ch == nil {
			L.ArgError(2, "Character expected")
			return 0
		}
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
