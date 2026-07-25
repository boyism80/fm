package entity

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
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			L.Push(lua.LString(group.Name))
			return 1
		},
		"set_property": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			key := L.CheckString(2)
			value := L.CheckString(3)
			group.SetProperty(key, value)
			return 0
		},
		"get_property": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			key := L.CheckString(2)
			L.Push(lua.LString(group.GetProperty(key)))
			return 1
		},
		"declare_maps": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			ids := make([]uint32, 0)
			if table, ok := L.Get(2).(*lua.LTable); ok {
				table.ForEach(func(_ lua.LValue, value lua.LValue) {
					if value.Type() == lua.LTNumber {
						ids = append(ids, uint32(lua.LVAsNumber(value)))
					}
				})
			} else {
				for i := 2; i <= L.GetTop(); i++ {
					ids = append(ids, uint32(L.CheckInt(i)))
				}
			}
			group.DeclareMaps(ids)
			return 0
		},
		"declare_min_players": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			group.DeclareMinPlayers(L.CheckInt(2))
			return 0
		},
		"declare_exit_map": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			group.DeclareExitMap(uint32(L.CheckInt(2)))
			return 0
		},
		"map": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			mapID := uint32(L.CheckInt(2))
			m := group.GetMap(mapID)
			if m == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, m))
			return 1
		},
		"get": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			id := L.CheckString(2)
			sm := group.Get(id)
			if sm == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, sm))
			return 1
		},
		"create": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			id := L.CheckString(2)
			opts := CreateOpts{}
			if L.GetTop() >= 3 && L.Get(3).Type() != lua.LTNil {
				leaderUd := L.CheckUserData(3)
				leader, ok := leaderUd.Value.(*Character)
				if !ok || leader == nil {
					L.ArgError(3, "Character expected")
					return 0
				}
				opts.Leader = leader
				opts.ScaleLevel = int(leader.GetLevel())
			}
			if L.GetTop() >= 4 && L.Get(4).Type() != lua.LTNil {
				partyUd := L.CheckUserData(4)
				party, ok := partyUd.Value.(*Party)
				if !ok || party == nil {
					L.ArgError(4, "Party expected")
					return 0
				}
				opts.Party = party
			}
			if L.GetTop() >= 5 {
				maxLevel := L.CheckInt(5)
				if maxLevel > 0 && opts.ScaleLevel > maxLevel {
					opts.ScaleLevel = maxLevel
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
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			leaderUd := L.CheckUserData(2)
			leader, ok := leaderUd.Value.(*Character)
			if !ok || leader == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			var party *Party
			if L.GetTop() >= 3 && L.Get(3).Type() != lua.LTNil {
				partyUd := L.CheckUserData(3)
				party, ok = partyUd.Value.(*Party)
				if !ok || party == nil {
					L.ArgError(3, "Party expected")
					return 0
				}
			} else {
				pid := leader.GetPartyID()
				if pid == nil || leader.GameWorld == nil {
					L.Push(lua.LNil)
					L.Push(lua.LString("no party"))
					return 2
				}
				party = leader.GameWorld.GetPartySystem().Get(*pid)
			}
			maxLevel := 0
			if L.GetTop() >= 4 {
				maxLevel = L.CheckInt(4)
			}
			sm, err := group.StartParty(leader, party, StartPartyOpts{MaxLevel: maxLevel})
			if err != nil {
				L.Push(lua.LNil)
				L.Push(lua.LString(err.Error()))
				return 2
			}
			L.Push(luax.NewLuable(L, sm))
			return 1
		},
		"start_solo": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			group, ok := ud.Value.(*StateMachineGroup)
			if !ok || group == nil {
				L.ArgError(1, "StateMachineGroup expected")
				return 0
			}
			playerUd := L.CheckUserData(2)
			player, ok := playerUd.Value.(*Character)
			if !ok || player == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			maxLevel := 0
			if L.GetTop() >= 3 {
				maxLevel = L.CheckInt(3)
			}
			sm, err := group.StartSolo(player, StartPartyOpts{MaxLevel: maxLevel})
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
