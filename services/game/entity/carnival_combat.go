package entity

import (
	"math/rand"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

func (m *CarnivalMatch) HandleTab(ch *Character, tab uint8, num int32) bool {
	if m == nil || ch == nil {
		return false
	}
	if m.GetState() != CarnivalStateBattle {
		return false
	}
	team := m.TeamOf(ch)
	if team == nil {
		return false
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil || mapInstance.Wz == nil || mapInstance.Wz.MonsterCarnival == nil {
		return false
	}
	mc := mapInstance.Wz.MonsterCarnival
	gw := m.gameWorld()
	if gw == nil {
		return false
	}
	resources := gw.GetResources()
	if resources == nil {
		return false
	}

	switch tab {
	case 0:
		idx := int(num)
		if idx < 0 || idx >= len(mc.Mobs) {
			return false
		}
		entry := mc.Mobs[idx]
		if entry.SpendCP <= 0 || !team.UseCP(ch, entry.SpendCP) {
			return false
		}
		pos, ok := pickCarnivalGenPos(mc.MobGenPos, team.TeamID, mapInstance, true)
		if !ok {
			return false
		}
		spawned, err := mapInstance.SpawnMob(entry.ID, pos, nil, constant.MobSpawnTypeNone, 0)
		if err != nil || spawned == nil {
			return false
		}
		spawned.CarnivalTeam = team.TeamID
		ch.Listener.OnCarnivalSummon(ch, tab, uint8(idx), ch.GetName())
		return true
	case 1:
		idx := int(num)
		if idx < 0 || idx >= len(mc.Skills) {
			return false
		}
		entry := mc.Skills[idx]
		skillDef := resources.GetMCSkill(entry.ID)
		if skillDef == nil {
			return false
		}
		levelData := skillDef.LevelData(1)
		if levelData == nil {
			return false
		}
		spend := entry.SpendCP
		if spend <= 0 {
			spend = levelData.CP
		}
		if spend <= 0 || !team.UseCP(ch, spend) {
			return false
		}
		enemy := m.enemyTeam(team.TeamID)
		if enemy != nil {
			for _, target := range enemy.Members(gw) {
				m.giveDebuff(resources, target, levelData.MobSkillID, levelData.MobSkillLevel)
			}
		}
		ch.Listener.OnCarnivalSummon(ch, tab, uint8(idx), ch.GetName())
		return true
	case 2:
		guardianDef := resources.GetMCGuardian(uint32(num))
		if guardianDef == nil {
			return false
		}
		levelData := guardianDef.LevelData(1)
		if levelData == nil {
			return false
		}
		if levelData.CP <= 0 || !team.UseCP(ch, levelData.CP) {
			return false
		}
		pos, ok := pickCarnivalGenPos(mc.GuardianGenPos, team.TeamID, mapInstance, false)
		if !ok {
			return false
		}
		reactorID := levelData.ReactorID
		if reactorID == 0 {
			reactorID = 9980000 + uint32(team.TeamID)
		}
		if _, err := mapInstance.SpawnReactorAt(reactorID, pos); err != nil {
			return false
		}
		ch.Listener.OnCarnivalSummon(ch, tab, uint8(num), ch.GetName())
		return true
	default:
		return false
	}
}

func (m *CarnivalMatch) enemyTeam(teamID constant.CarnivalTeam) *CarnivalTeam {
	if teamID == constant.CarnivalTeamRed {
		return m.Team(constant.CarnivalTeamBlue)
	}
	return m.Team(constant.CarnivalTeamRed)
}

func pickCarnivalGenPos(gens []wz.CarnivalGenPos, teamID constant.CarnivalTeam, mapInstance *Map, checkMobOccupancy bool) (types.Point[int16], bool) {
	if len(gens) == 0 || mapInstance == nil {
		return types.Point[int16]{}, false
	}
	candidates := make([]types.Point[int16], 0, len(gens))
	for _, gen := range gens {
		if gen.Team != -1 && gen.Team != int(teamID) {
			continue
		}
		if checkMobOccupancy && carnivalGenPosOccupied(mapInstance, gen.Pos, teamID) {
			continue
		}
		candidates = append(candidates, gen.Pos)
	}
	if len(candidates) == 0 {
		return types.Point[int16]{}, false
	}
	return candidates[rand.Intn(len(candidates))], true
}

func carnivalGenPosOccupied(mapInstance *Map, pos types.Point[int16], teamID constant.CarnivalTeam) bool {
	for _, obj := range mapInstance.GetMobs() {
		mob, ok := obj.(*Mob)
		if !ok || mob == nil {
			continue
		}
		if mob.CarnivalTeam != teamID && mob.CarnivalTeam != constant.CarnivalTeamNone {
			continue
		}
		p := mob.GetPosition()
		if p.X == pos.X && (p.Y == pos.Y || p.Y == pos.Y-1) {
			return true
		}
	}
	return false
}

func (m *CarnivalMatch) giveDebuff(resources *wz.Resources, target *Character, skillID uint32, skillLevel uint8) {
	if target == nil || resources == nil || skillID == 0 {
		return
	}
	mobSkill := resources.GetMobSkill(skillID, skillLevel)
	if mobSkill == nil {
		return
	}
	for _, flag := range constant.AllDebuffFlags() {
		if flag.DiseaseSkillID == uint16(skillID) {
			duration := time.Duration(mobSkill.DurationMs) * time.Millisecond
			if duration <= 0 {
				duration = 5 * time.Second
			}
			target.GiveDebuff(flag, duration, int16(mobSkill.X), uint16(skillID), uint16(skillLevel))
			return
		}
	}
}

func (ch *Character) PickupCarnivalItem(wzConsume *wz.Consume) bool {
	if ch == nil || wzConsume == nil {
		return false
	}
	match := ch.CarnivalMatch()
	if match == nil {
		return false
	}
	team := match.TeamOf(ch)
	if team == nil {
		return false
	}
	if wzConsume.CP > 0 {
		team.AddCP(ch, wzConsume.CP)
	}
	if wzConsume.NuffSkillID > 0 {
		enemy := match.enemyTeam(team.TeamID)
		if enemy != nil {
			resources := ch.GameWorld.GetResources()
			for _, target := range enemy.Members(ch.GameWorld) {
				match.giveDebuff(resources, target, wzConsume.NuffSkillID, wzConsume.NuffSkillLevel)
			}
		}
	}
	return true
}
