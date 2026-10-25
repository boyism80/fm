package game

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	lua "github.com/yuin/gopher-lua"
)

func registerCarnivalConstants(L *lua.LState) {
	teamTable := L.NewTable()
	teamTable.RawSetString("NONE", lua.LNumber(constant.CarnivalTeamNone))
	teamTable.RawSetString("RED", lua.LNumber(constant.CarnivalTeamRed))
	teamTable.RawSetString("BLUE", lua.LNumber(constant.CarnivalTeamBlue))
	L.SetGlobal("CARNIVAL_TEAM", teamTable)

	stateTable := L.NewTable()
	stateTable.RawSetString("EMPTY", lua.LNumber(CarnivalStateEmpty))
	stateTable.RawSetString("WAITING", lua.LNumber(CarnivalStateWaiting))
	stateTable.RawSetString("READY", lua.LNumber(CarnivalStateReady))
	stateTable.RawSetString("BATTLE", lua.LNumber(CarnivalStateBattle))
	stateTable.RawSetString("REWARD", lua.LNumber(CarnivalStateReward))
	L.SetGlobal("CARNIVAL_STATE", stateTable)

	resultTable := L.NewTable()
	resultTable.RawSetString("UNDECIDED", lua.LNumber(CarnivalResultUndecided))
	resultTable.RawSetString("RED_WIN", lua.LNumber(CarnivalResultRedWin))
	resultTable.RawSetString("BLUE_WIN", lua.LNumber(CarnivalResultBlueWin))
	resultTable.RawSetString("DRAW", lua.LNumber(CarnivalResultDraw))
	L.SetGlobal("CARNIVAL_RESULT", resultTable)

	tabTable := L.NewTable()
	tabTable.RawSetString("MOB", lua.LNumber(constant.CarnivalTabMob))
	tabTable.RawSetString("SKILL", lua.LNumber(constant.CarnivalTabSkill))
	tabTable.RawSetString("GUARDIAN", lua.LNumber(constant.CarnivalTabGuardian))
	L.SetGlobal("CARNIVAL_TAB", tabTable)

	summonTable := L.NewTable()
	summonTable.RawSetString("SUCCESS", lua.LNumber(CarnivalSummonSuccess))
	summonTable.RawSetString("LACK_CP", lua.LNumber(CarnivalSummonLackCP))
	summonTable.RawSetString("NO_SLOT", lua.LNumber(CarnivalSummonNoSlot))
	summonTable.RawSetString("FAILED", lua.LNumber(CarnivalSummonFailed))
	L.SetGlobal("CARNIVAL_SUMMON_RESULT", summonTable)
}

func (m *CarnivalMatch) LuaTypeName() string { return "LuaCarnivalMatch" }

func (m *CarnivalMatch) String() string { return m.LuaTypeName() }

func (m *CarnivalMatch) Type() lua.LValueType { return lua.LTUserData }

func (m *CarnivalMatch) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.waitingMapID))
			return 1
		},
		"slot": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.slotIndex))
			return 1
		},
		"state": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.GetState()))
			return 1
		},
		"max_members": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.maxMembers))
			return 1
		},
		"waiting_map_id": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.waitingMapID))
			return 1
		},
		"field_map_id": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.fieldMapID))
			return 1
		},
		"revive_map_id": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.reviveMapID))
			return 1
		},
		"win_map_id": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.winMapID))
			return 1
		},
		"lose_map_id": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.loseMapID))
			return 1
		},
		"red_team": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.Push(lua.LNil)
				return 1
			}
			team := match.Team(constant.CarnivalTeamRed)
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"blue_team": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.Push(lua.LNil)
				return 1
			}
			team := match.Team(constant.CarnivalTeamBlue)
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"enemy_team": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			team := match.enemyTeam(constant.CarnivalTeam(L.CheckInt(2)))
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"find_team": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			ch := luax.Check[Character](L, 2)
			team := match.FindTeam(ch)
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"has_pending_challenge": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LBool(match.HasPendingChallenge()))
			return 1
		},
		"pending_challenge_members": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			members := match.PendingChallengeMembers()
			tbl := L.NewTable()
			for i, member := range members {
				tbl.RawSetInt(i+1, luax.NewLuable(L, member))
			}
			L.Push(tbl)
			L.Push(lua.LNumber(len(members)))
			return 2
		},
		"accept_pending_challenge": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			gw := match.gameWorld()
			L.Push(lua.LBool(match.AcceptPendingChallenge(gw)))
			return 1
		},
		"reject_pending_challenge": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			match.RejectPendingChallenge()
			return 0
		},
		"result": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			L.Push(lua.LNumber(match.Result()))
			return 1
		},
		"finish": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			gw := match.gameWorld()
			L.Push(lua.LBool(match.Finish(gw)))
			return 1
		},
		"conclude": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			gw := match.gameWorld()
			L.Push(lua.LBool(match.Conclude(gw)))
			return 1
		},
		"set_state": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			match.SetState(CarnivalState(L.CheckNumber(2)))
			L.Push(lua.LTrue)
			return 1
		},
		"show_result": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			match.ShowResult()
			return 0
		},
		"on_mob_spawn": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			mob := luax.Check[Mob](L, 2)
			match.OnMobSpawn(mob)
			return 0
		},
		"destroy_guardian": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			reactor := luax.Check[Reactor](L, 2)
			match.DestroyGuardian(reactor)
			return 0
		},
		"on_player_died": func(L *lua.LState) int {
			match := luax.Check[CarnivalMatch](L, 1)
			ch := luax.Check[Character](L, 2)
			match.OnPlayerDied(ch)
			L.Push(lua.LTrue)
			return 1
		},
	}
}

func (t *CarnivalTeam) LuaTypeName() string { return "LuaCarnivalTeam" }

func (t *CarnivalTeam) String() string { return t.LuaTypeName() }

func (t *CarnivalTeam) Type() lua.LValueType { return lua.LTUserData }

func (t *CarnivalTeam) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			L.Push(lua.LNumber(team.teamID))
			return 1
		},
		"leader": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil || team.match == nil {
				L.Push(lua.LNil)
				return 1
			}
			leader := team.Leader()
			if leader == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, leader))
			}
			return 1
		},
		"members": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil || team.match == nil {
				L.Push(lua.LNil)
				return 1
			}
			tbl := L.NewTable()
			for i, ch := range team.Members() {
				tbl.RawSetInt(i+1, luax.NewLuable(L, ch))
			}
			L.Push(tbl)
			return 1
		},
		"roster": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			tbl := L.NewTable()
			for i, member := range team.roster {
				tbl.RawSetInt(i+1, luax.NewLuable(L, member))
			}
			L.Push(tbl)
			return 1
		},
		"size": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			L.Push(lua.LNumber(len(team.roster)))
			return 1
		},
		"available_cp": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			L.Push(lua.LNumber(team.availableCP))
			return 1
		},
		"total_cp": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			L.Push(lua.LNumber(team.totalCP))
			return 1
		},
		"personal_cp": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			ch := luax.Check[Character](L, 2)
			personal := team.PersonalCP(ch.GetID())
			if personal == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, personal))
			}
			return 1
		},
		"add_cp": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			ch := luax.Check[Character](L, 2)
			amount := int(L.CheckNumber(3))
			L.Push(lua.LBool(team.AddCP(ch, amount)))
			return 1
		},
		"summon_mob": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			ch := luax.Check[Character](L, 2)
			field := luax.Check[Map](L, 3)
			L.Push(lua.LNumber(team.SummonMob(ch, field, L.CheckInt(4))))
			return 1
		},
		"use_skill": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			ch := luax.Check[Character](L, 2)
			field := luax.Check[Map](L, 3)
			L.Push(lua.LNumber(team.UseSkill(ch, field, L.CheckInt(4))))
			return 1
		},
		"summon_guardian": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			ch := luax.Check[Character](L, 2)
			field := luax.Check[Map](L, 3)
			L.Push(lua.LNumber(team.SummonGuardian(ch, field, uint32(L.CheckInt(4)))))
			return 1
		},
		"is_winner": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			L.Push(lua.LBool(team.winner))
			return 1
		},
		"set_winner": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			team.winner = lua.LVAsBool(L.Get(2))
			L.Push(lua.LTrue)
			return 1
		},
		"warp": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			mapID := uint32(L.CheckNumber(2))
			portalName := L.OptString(3, "")
			cfg, _ := luax.GetConfiguration(L)
			L.Push(lua.LBool(team.Warp(cfg.ActorContext, mapID, portalName)))
			return 1
		},
		"remove_member": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			ch := luax.Check[Character](L, 2)
			team.RemoveMember(ch)
			return 0
		},
		"debuff": func(L *lua.LState) int {
			team := luax.Check[CarnivalTeam](L, 1)
			field := luax.Check[Map](L, 2)
			L.Push(lua.LBool(team.Debuff(field, uint32(L.CheckNumber(3)))))
			return 1
		},
	}
}

func (p *CarnivalPersonalCP) LuaTypeName() string { return "LuaCarnivalPersonalCP" }

func (p *CarnivalPersonalCP) String() string { return p.LuaTypeName() }

func (p *CarnivalPersonalCP) Type() lua.LValueType { return lua.LTUserData }

func (p *CarnivalPersonalCP) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"available_cp": func(L *lua.LState) int {
			personal := luax.Check[CarnivalPersonalCP](L, 1)
			L.Push(lua.LNumber(personal.availableCP))
			return 1
		},
		"total_cp": func(L *lua.LState) int {
			personal := luax.Check[CarnivalPersonalCP](L, 1)
			L.Push(lua.LNumber(personal.totalCP))
			return 1
		},
	}
}

func (member *CarnivalMember) LuaTypeName() string { return "LuaCarnivalMember" }

func (member *CarnivalMember) String() string { return member.LuaTypeName() }

func (member *CarnivalMember) Type() lua.LValueType { return lua.LTUserData }

func (member *CarnivalMember) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			L.Push(lua.LNumber(luax.Check[CarnivalMember](L, 1).id))
			return 1
		},
		"name": func(L *lua.LState) int {
			L.Push(lua.LString(luax.Check[CarnivalMember](L, 1).name))
			return 1
		},
		"level": func(L *lua.LState) int {
			L.Push(lua.LNumber(luax.Check[CarnivalMember](L, 1).level))
			return 1
		},
		"class": func(L *lua.LState) int {
			L.Push(lua.LNumber(luax.Check[CarnivalMember](L, 1).class))
			return 1
		},
	}
}

func RegisterCarnivalLua(gw GameWorld, L *lua.LState) {
	registerCarnivalConstants(L)
	luax.RegisterLuaType[*CarnivalMatch](L)
	luax.RegisterLuaType[*CarnivalTeam](L)
	luax.RegisterLuaType[*CarnivalPersonalCP](L)
	luax.RegisterLuaType[*CarnivalMember](L)

	carnivalTable := L.NewTable()
	reg := gw.Systems().Get[*CarnivalSystem]()

	carnivalTable.RawSetString("bind_group", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			return 0
		}
		group := luax.Check[StateMachineGroup](L, 1)
		reg.BindGroup(group)
		return 0
	}))
	carnivalTable.RawSetString("group", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			L.Push(lua.LNil)
			return 1
		}
		group := reg.Group()
		if group == nil {
			L.Push(lua.LNil)
		} else {
			L.Push(luax.NewLuable(L, group))
		}
		return 1
	}))
	carnivalTable.RawSetString("register", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			L.Push(lua.LNil)
			return 1
		}
		slot := int(L.CheckNumber(1))
		waitingMapID := uint32(L.CheckNumber(2))
		fieldMapID := uint32(L.CheckNumber(3))
		reviveMapID := uint32(L.CheckNumber(4))
		winMapID := uint32(L.CheckNumber(5))
		loseMapID := uint32(L.CheckNumber(6))
		maxMembers := int(L.CheckNumber(7))
		match := reg.Register(slot, waitingMapID, fieldMapID, reviveMapID, winMapID, loseMapID, maxMembers)
		L.Push(luax.NewLuable(L, match))
		return 1
	}))
	carnivalTable.RawSetString("match", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			L.Push(lua.LNil)
			return 1
		}
		slot := int(L.CheckNumber(1))
		match := reg.Match(slot)
		if match == nil {
			L.Push(lua.LNil)
		} else {
			L.Push(luax.NewLuable(L, match))
		}
		return 1
	}))
	carnivalTable.RawSetString("matches", L.NewFunction(func(L *lua.LState) int {
		tbl := L.NewTable()
		if reg != nil {
			i := 1
			for _, match := range reg.Matches() {
				tbl.RawSetInt(i, luax.NewLuable(L, match))
				i++
			}
		}
		L.Push(tbl)
		return 1
	}))
	carnivalTable.RawSetString("map_match", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			L.Push(lua.LNil)
			return 1
		}
		mapID := uint32(L.CheckNumber(1))
		match := reg.MapMatch(mapID)
		if match == nil {
			L.Push(lua.LNil)
		} else {
			L.Push(luax.NewLuable(L, match))
		}
		return 1
	}))
	carnivalTable.RawSetString("set_skill_hit_chance", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			return 0
		}
		reg.SetSkillHitChance(uint32(L.CheckNumber(1)), int(L.CheckNumber(2)))
		return 0
	}))
	carnivalTable.RawSetString("enter", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			L.Push(lua.LFalse)
			return 1
		}
		slot := int(L.CheckNumber(1))
		ch := luax.Check[Character](L, 2)
		L.Push(lua.LBool(reg.Enter(slot, ch)))
		return 1
	}))
	carnivalTable.RawSetString("challenge", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			L.Push(lua.LFalse)
			return 1
		}
		slot := int(L.CheckNumber(1))
		ch := luax.Check[Character](L, 2)
		L.Push(lua.LBool(reg.Challenge(slot, ch)))
		return 1
	}))
	L.SetGlobal("carnival", carnivalTable)
}
