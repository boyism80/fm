package entity

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
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
}

func (m *CarnivalMatch) LuaTypeName() string { return "LuaCarnivalMatch" }

func (m *CarnivalMatch) String() string { return m.LuaTypeName() }

func (m *CarnivalMatch) Type() lua.LValueType { return lua.LTUserData }

func (m *CarnivalMatch) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.WaitingMapID))
			return 1
		},
		"slot": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.SlotIndex))
			return 1
		},
		"state": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.GetState()))
			return 1
		},
		"max_members": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.MaxMembers))
			return 1
		},
		"waiting_map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.WaitingMapID))
			return 1
		},
		"field_map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.FieldMapID))
			return 1
		},
		"revive_map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.ReviveMapID))
			return 1
		},
		"win_map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.WinMapID))
			return 1
		},
		"lose_map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.LoseMapID))
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
		"team": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			teamID := constant.CarnivalTeam(L.CheckInt(2))
			team := match.Team(teamID)
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"enemy_team": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			team := match.enemyTeam(constant.CarnivalTeam(L.CheckInt(2)))
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"find_team": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			chUD := L.CheckUserData(2)
			ch, ok := chUD.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			team := match.FindTeam(ch)
			if team == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, team))
			}
			return 1
		},
		"has_pending_challenge": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LBool(match.HasPendingChallenge()))
			return 1
		},
		"pending_challenge_members": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			members, size := match.PendingChallengeMembers(match.gameWorld())
			tbl := L.NewTable()
			for i, ch := range members {
				tbl.RawSetInt(i+1, luax.NewLuable(L, ch))
			}
			L.Push(tbl)
			L.Push(lua.LNumber(size))
			return 2
		},
		"accept_pending_challenge": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			gw := match.gameWorld()
			L.Push(lua.LBool(match.AcceptPendingChallenge(gw)))
			return 1
		},
		"reject_pending_challenge": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			_, shouldOpen := match.RejectPendingChallenge()
			L.Push(lua.LBool(shouldOpen))
			return 1
		},
		"result": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			L.Push(lua.LNumber(match.Result()))
			return 1
		},
		"warp_all": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			mapID := uint32(L.CheckNumber(2))
			portalName := L.OptString(3, "")
			cfg, _ := luax.GetConfiguration(L)
			L.Push(lua.LBool(match.WarpAll(cfg.ActorContext, mapID, portalName)))
			return 1
		},
		"finish": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			gw := match.gameWorld()
			L.Push(lua.LBool(match.Finish(gw)))
			return 1
		},
		"conclude": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			gw := match.gameWorld()
			L.Push(lua.LBool(match.Conclude(gw)))
			return 1
		},
		"set_state": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			match.SetState(CarnivalState(L.CheckNumber(2)))
			L.Push(lua.LTrue)
			return 1
		},
		"on_player_died": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			match, ok := ud.Value.(*CarnivalMatch)
			if !ok || match == nil {
				L.ArgError(1, "CarnivalMatch expected")
				return 0
			}
			chUD := L.CheckUserData(2)
			ch, ok := chUD.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
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
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			L.Push(lua.LNumber(team.TeamID))
			return 1
		},
		"leader": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil || team.Match == nil {
				L.Push(lua.LNil)
				return 1
			}
			gw := team.Match.gameWorld()
			leader := team.Leader(gw)
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
			if !ok || team == nil || team.Match == nil {
				L.Push(lua.LNil)
				return 1
			}
			gw := team.Match.gameWorld()
			tbl := L.NewTable()
			for i, ch := range team.Members(gw) {
				tbl.RawSetInt(i+1, luax.NewLuable(L, ch))
			}
			L.Push(tbl)
			return 1
		},
		"size": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			L.Push(lua.LNumber(len(team.MemberIDs)))
			return 1
		},
		"available_cp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			L.Push(lua.LNumber(team.AvailableCP))
			return 1
		},
		"total_cp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			L.Push(lua.LNumber(team.TotalCP))
			return 1
		},
		"personal_cp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			chUD := L.CheckUserData(2)
			ch, ok := chUD.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			personal := team.PersonalCP(ch.GetID())
			if personal == nil {
				L.Push(lua.LNil)
			} else {
				L.Push(luax.NewLuable(L, personal))
			}
			return 1
		},
		"add_cp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			chUD := L.CheckUserData(2)
			ch, ok := chUD.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			amount := int(L.CheckNumber(3))
			L.Push(lua.LBool(team.AddCP(ch, amount)))
			return 1
		},
		"use_cp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			chUD := L.CheckUserData(2)
			ch, ok := chUD.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			amount := int(L.CheckNumber(3))
			L.Push(lua.LBool(team.UseCP(ch, amount)))
			return 1
		},
		"is_winner": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			L.Push(lua.LBool(team.Winner))
			return 1
		},
		"set_winner": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			team.Winner = lua.LVAsBool(L.Get(2))
			L.Push(lua.LTrue)
			return 1
		},
		"warp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			mapID := uint32(L.CheckNumber(2))
			portalName := L.OptString(3, "")
			cfg, _ := luax.GetConfiguration(L)
			L.Push(lua.LBool(team.Warp(cfg.ActorContext, mapID, portalName)))
			return 1
		},
		"clear": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil || team.Match == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			gw := team.Match.gameWorld()
			team.Clear(gw)
			L.Push(lua.LTrue)
			return 1
		},
		"remove_member": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			chUD := L.CheckUserData(2)
			ch, ok := chUD.Value.(*Character)
			if !ok || ch == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			team.RemoveMember(ch)
			return 0
		},
		"debuff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			team, ok := ud.Value.(*CarnivalTeam)
			if !ok || team == nil {
				L.ArgError(1, "CarnivalTeam expected")
				return 0
			}
			mapUD := L.CheckUserData(2)
			field, ok := mapUD.Value.(*Map)
			if !ok || field == nil {
				L.ArgError(2, "Map expected")
				return 0
			}
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
			ud := L.CheckUserData(1)
			personal, ok := ud.Value.(*CarnivalPersonalCP)
			if !ok || personal == nil {
				L.ArgError(1, "CarnivalPersonalCP expected")
				return 0
			}
			L.Push(lua.LNumber(personal.AvailableCP))
			return 1
		},
		"total_cp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			personal, ok := ud.Value.(*CarnivalPersonalCP)
			if !ok || personal == nil {
				L.ArgError(1, "CarnivalPersonalCP expected")
				return 0
			}
			L.Push(lua.LNumber(personal.TotalCP))
			return 1
		},
	}
}

func RegisterCarnivalLua(L *lua.LState, gw GameWorld) {
	registerCarnivalConstants(L)
	luax.RegisterLuaType[*CarnivalMatch](L)
	luax.RegisterLuaType[*CarnivalTeam](L)
	luax.RegisterLuaType[*CarnivalPersonalCP](L)

	carnivalTable := L.NewTable()
	reg := gw.GetCarnivalRegistry()

	carnivalTable.RawSetString("bind_group", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			return 0
		}
		ud := L.CheckUserData(1)
		group, ok := ud.Value.(*StateMachineGroup)
		if !ok || group == nil {
			L.ArgError(1, "StateMachineGroup expected")
			return 0
		}
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
	carnivalTable.RawSetString("skill", L.NewFunction(func(L *lua.LState) int {
		skill := gw.GetResources().GetCarnivalSkill(uint32(L.CheckNumber(1)))
		if skill == nil {
			L.Push(lua.LNil)
			return 1
		}
		tbl := L.NewTable()
		tbl.RawSetString("spend_cp", lua.LNumber(skill.SpendCP))
		tbl.RawSetString("targets_all", lua.LBool(skill.TargetsAll))
		L.Push(tbl)
		return 1
	}))
	carnivalTable.RawSetString("guardian", L.NewFunction(func(L *lua.LState) int {
		resources := gw.GetResources()
		guardian := resources.GetCarnivalGuardian(uint32(L.CheckNumber(1)))
		if guardian == nil {
			L.Push(lua.LNil)
			return 1
		}
		tbl := L.NewTable()
		tbl.RawSetString("spend_cp", lua.LNumber(guardian.SpendCP))
		tbl.RawSetString("mob_skill_id", lua.LNumber(guardian.MobSkillID))
		if levelData := resources.GetMobSkill(guardian.MobSkillID, guardian.Level); levelData != nil {
			tbl.RawSetString("skill", luax.NewLuable(L, &MobSkill{
				Slot:      wz.MobSkillSlot{SkillID: guardian.MobSkillID, Level: guardian.Level},
				LevelData: levelData,
			}))
		}
		L.Push(tbl)
		return 1
	}))
	carnivalTable.RawSetString("enter", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			L.Push(lua.LFalse)
			return 1
		}
		slot := int(L.CheckNumber(1))
		chUD := L.CheckUserData(2)
		ch, ok := chUD.Value.(*Character)
		if !ok || ch == nil {
			L.ArgError(2, "Character expected")
			return 0
		}
		L.Push(lua.LBool(reg.Enter(slot, ch)))
		return 1
	}))
	carnivalTable.RawSetString("challenge", L.NewFunction(func(L *lua.LState) int {
		if reg == nil {
			L.Push(lua.LFalse)
			L.Push(lua.LFalse)
			return 2
		}
		slot := int(L.CheckNumber(1))
		chUD := L.CheckUserData(2)
		ch, ok := chUD.Value.(*Character)
		if !ok || ch == nil {
			L.ArgError(2, "Character expected")
			return 0
		}
		succeeded, shouldOpen := reg.Challenge(slot, ch)
		L.Push(lua.LBool(succeeded))
		L.Push(lua.LBool(shouldOpen))
		return 2
	}))
	L.SetGlobal("carnival", carnivalTable)
}
