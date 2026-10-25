package game

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (g *StateMachineGroup) LuaTypeName() string {
	return "LuaStateMachineGroup"
}

func (g *StateMachineGroup) String() string {
	return g.LuaTypeName()
}

func (g *StateMachineGroup) Type() lua.LValueType {
	return lua.LTUserData
}

func (g *StateMachineGroup) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"full": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			L.Push(lua.LBool(group.Full()))
			return 1
		},
		"name": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			L.Push(lua.LString(group.name))
			return 1
		},
		"set_property": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			key := L.CheckString(2)
			value := L.CheckString(3)
			group.SetProperty(key, value)
			return 0
		},
		"get_property": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			key := L.CheckString(2)
			L.Push(lua.LString(group.GetProperty(key)))
			return 1
		},
		"min_players": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			group.SetMinPlayers(L.CheckInt(2))
			return 0
		},
		"exit_map": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			group.SetExitMap(uint32(L.CheckInt(2)))
			return 0
		},
		"max_machines": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			if L.Get(2) == lua.LNil {
				group.RemoveMachineLimit()
				return 0
			}
			group.SetMaxMachines(L.CheckInt(2))
			return 0
		},
		"override_max_machines": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			if L.Get(2) == lua.LNil {
				group.OverrideMachineCap(MachineCap{})
				return 0
			}
			group.OverrideMachineCap(MachineCap{limited: true, max: L.CheckInt(2)})
			return 0
		},
		"reset_max_machines": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			group.ClearMachineCapOverride()
			return 0
		},
		"machine_cap": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			c, overridden := group.MachineCap()
			if c.limited {
				L.Push(lua.LNumber(c.max))
			} else {
				L.Push(lua.LNil)
			}
			L.Push(lua.LBool(overridden))
			return 2
		},
		"get": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			id := L.CheckString(2)
			sm := group.Get(id)
			if sm == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, sm))
			return 1
		},
		"machines": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			tbl := L.NewTable()
			for _, sm := range group.Machines() {
				tbl.Append(luax.NewLuable(L, sm))
			}
			L.Push(tbl)
			return 1
		},
		"create": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			id := L.CheckString(2)
			opts := CreateOpts{}
			if L.GetTop() >= 3 && L.Get(3).Type() != lua.LTNil {
				leader := luax.Check[Character](L, 3)
				opts.leader = leader
				opts.scaleLevel = int(leader.GetLevel())
			}
			if L.GetTop() >= 4 && L.Get(4).Type() != lua.LTNil {
				party := luax.Check[Party](L, 4)
				opts.party = party
			}
			if L.GetTop() >= 5 {
				maxLevel := L.CheckInt(5)
				if maxLevel > 0 && opts.scaleLevel > maxLevel {
					opts.scaleLevel = maxLevel
				}
			}
			sm, err := group.Create(id, opts)
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			L.Push(luax.NewLuable(L, sm))
			return 1
		},
		"start_party": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			leader := luax.Check[Character](L, 2)
			var party *Party
			if L.GetTop() >= 3 && L.Get(3).Type() != lua.LTNil {
				party = luax.Check[Party](L, 3)
			} else {
				pid := leader.party.ID()
				if pid == nil || leader.gameWorld == nil {
					L.Push(lua.LNil)
					L.Push(lua.LString("no party"))
					return 2
				}
				party = leader.gameWorld.Systems().Get[*PartySystem]().Get(*pid)
			}
			maxLevel := 0
			if L.GetTop() >= 4 && L.Get(4).Type() == lua.LTNumber {
				maxLevel = L.CheckInt(4)
			}
			sm, err := group.StartParty(leader, party, StartPartyOpts{maxLevel: maxLevel, props: ParseMachineProps(L)})
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			L.Push(luax.NewLuable(L, sm))
			return 1
		},
		"start_solo": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			player := luax.Check[Character](L, 2)
			maxLevel := 0
			if L.GetTop() >= 3 && L.Get(3).Type() == lua.LTNumber {
				maxLevel = L.CheckInt(3)
			}
			sm, err := group.StartSolo(player, StartPartyOpts{maxLevel: maxLevel, props: ParseMachineProps(L)})
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			L.Push(luax.NewLuable(L, sm))
			return 1
		},
		"start_persistent": func(L *lua.LState) int {
			group := luax.Check[StateMachineGroup](L, 1)
			id := "persistent"
			if L.GetTop() >= 2 && L.Get(2).Type() == lua.LTString {
				id = L.CheckString(2)
			}
			sm, err := group.StartPersistent(id, ParseMachineProps(L))
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			L.Push(luax.NewLuable(L, sm))
			return 1
		},
	}
}

// ParseMachineProps copies the first table among a start call's arguments into the machine's starting properties.
// The table belongs to the caller's Lua state, so only its string-convertible values cross over.
func ParseMachineProps(L *lua.LState) map[string]string {
	for i := 2; i <= L.GetTop(); i++ {
		tbl, ok := L.Get(i).(*lua.LTable)
		if !ok {
			continue
		}
		props := make(map[string]string)
		tbl.ForEach(func(k, v lua.LValue) {
			props[k.String()] = v.String()
		})
		return props
	}
	return nil
}
