package entity

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

func (m *CarnivalMatch) Summon(ch *Character, tab uint8, num int32) bool {
	if m.GetState() != CarnivalStateBattle {
		return false
	}
	team := m.FindTeam(ch)
	mapInstance := ch.GetMap()
	gw := m.gameWorld()
	if team == nil || mapInstance == nil || mapInstance.Wz == nil || mapInstance.Wz.Carnival == nil || gw == nil {
		return false
	}
	resources := gw.GetResources()
	if resources == nil {
		return false
	}
	mc := mapInstance.Wz.Carnival

	switch tab {
	case 0:
		idx := int(num)
		if idx >= 0 && idx < len(mc.Mobs) {
			entry := mc.Mobs[idx]
			pos, posOK := m.pickGenPos(mc.MobGenPos, team.TeamID, mapInstance, true)
			personal := team.PersonalCP(ch.GetID())
			if entry.SpendCP > 0 && posOK && personal != nil && personal.AvailableCP >= entry.SpendCP && team.AvailableCP >= entry.SpendCP {
				mob, err := mapInstance.SpawnMob(entry.ID, pos, nil, constant.MobSpawnTypeNone, 0)
				if err == nil && mob != nil && team.UseCP(ch, entry.SpendCP) {
					mob.CarnivalTeam = team.TeamID
					ch.Listener.OnCarnivalSummon(ch, tab, uint8(idx), ch.GetName())
					return true
				}
			}
		}
		return false
	case 1:
		idx := int(num)
		if idx >= 0 && idx < len(mc.Skills) {
			entry := mc.Skills[idx]
			skillDef := resources.GetCarnivalSkill(entry.ID)
			if skillDef != nil {
				spend := entry.SpendCP
				if spend <= 0 {
					spend = skillDef.SpendCP
				}
				enemy := m.enemyTeam(team.TeamID)
				personal := team.PersonalCP(ch.GetID())
				if spend > 0 && enemy != nil && personal != nil && personal.AvailableCP >= spend && team.AvailableCP >= spend && m.debuffEnemies(gw, resources, enemy, entry.ID, skillDef) && team.UseCP(ch, spend) {
					ch.Listener.OnCarnivalSummon(ch, tab, uint8(idx), ch.GetName())
					return true
				}
			}
		}
		return false
	case 2:
		guardianDef := resources.GetCarnivalGuardian(uint32(num))
		if guardianDef != nil && guardianDef.SpendCP > 0 {
			personal := team.PersonalCP(ch.GetID())
			reactorName := fmt.Sprintf("%d%d", team.TeamID, num)
			existing := mapInstance.FindReactorName(reactorName)
			pos, posOK := m.pickGenPos(mc.GuardianGenPos, team.TeamID, mapInstance, false)
			if personal != nil && personal.AvailableCP >= guardianDef.SpendCP && team.AvailableCP >= guardianDef.SpendCP && (existing == nil || existing.State >= 5) && posOK {
				reactorID := 9980000 + uint32(team.TeamID)
				reactor, err := mapInstance.SpawnReactorTemplate(reactorID, pos, reactorName)
				if err == nil && reactor != nil && team.UseCP(ch, guardianDef.SpendCP) {
					m.buffAllyMobs(mapInstance, team.TeamID, resources, guardianDef.MobSkillID, guardianDef.Level)
					ch.Listener.OnCarnivalSummon(ch, tab, uint8(num), ch.GetName())
					return true
				}
			}
		}
		return false
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

func (m *CarnivalMatch) debuffEnemies(gw GameWorld, resources *wz.Resources, enemy *CarnivalTeam, catalogID uint32, skillDef *wz.CarnivalSkill) bool {
	if enemy == nil || resources == nil || skillDef == nil {
		return false
	}
	targets := enemy.Members(gw)
	if len(targets) == 0 {
		return false
	}
	rand.Shuffle(len(targets), func(i, j int) {
		targets[i], targets[j] = targets[j], targets[i]
	})
	chance := m.registry.SkillHitChance(catalogID, skillDef.HitChance)
	if !skillDef.TargetsAll {
		targets = targets[:1]
	}
	applied := false
	for _, target := range targets {
		if skillDef.TargetsAll && rand.Intn(100) >= chance {
			continue
		}
		if !m.giveDebuff(resources, target, skillDef.MobSkillID, skillDef.Level) {
			target.Buffs.Dispel()
		}
		applied = true
	}
	return applied
}

func (m *CarnivalMatch) buffAllyMobs(mapInstance *Map, teamID constant.CarnivalTeam, resources *wz.Resources, skillID uint32, skillLevel uint8) {
	if mapInstance == nil || resources == nil || skillID == 0 {
		return
	}
	levelData := resources.GetMobSkill(skillID, skillLevel)
	if levelData == nil {
		return
	}
	flag, ok := m.mobSkillBuffFlag(skillID)
	if !ok {
		return
	}
	duration := time.Duration(levelData.DurationMs) * time.Millisecond
	if duration <= 0 {
		duration = 60 * time.Second
	}
	for _, obj := range mapInstance.GetMobs() {
		mob, ok := obj.(*Mob)
		if !ok || mob == nil || mob.CarnivalTeam != teamID {
			continue
		}
		mob.GiveMobBuff(flag, int32(levelData.X), duration, nil, skillLevel, 0, 1)
	}
}

func (m *CarnivalMatch) mobSkillBuffFlag(skillID uint32) (constant.MobBuffFlag, bool) {
	switch skillID {
	case 100, 110, 150:
		return constant.MobBuffWeaponAttackUp, true
	case 101, 111, 151:
		return constant.MobBuffMagicAttackUp, true
	case 102, 112, 152:
		return constant.MobBuffWeaponDefenseUp, true
	case 103, 113, 153:
		return constant.MobBuffMagicDefenseUp, true
	case 154:
		return constant.MobBuffAcc, true
	case 155:
		return constant.MobBuffAvoid, true
	case 115, 156:
		return constant.MobBuffSpeed, true
	case 140:
		return constant.MobBuffWeaponImmunity, true
	case 141:
		return constant.MobBuffMagicImmunity, true
	default:
		return 0, false
	}
}

func (m *CarnivalMatch) pickGenPos(gens []wz.CarnivalGenPos, teamID constant.CarnivalTeam, mapInstance *Map, checkMobOccupancy bool) (types.Point[int16], bool) {
	if len(gens) == 0 || mapInstance == nil {
		return types.Point[int16]{}, false
	}
	candidates := make([]types.Point[int16], 0, len(gens))
	for _, gen := range gens {
		if gen.Team != -1 && gen.Team != int(teamID) {
			continue
		}
		if checkMobOccupancy && m.genPosOccupied(mapInstance, gen.Pos, teamID) {
			continue
		}
		candidates = append(candidates, gen.Pos)
	}
	if len(candidates) == 0 {
		return types.Point[int16]{}, false
	}
	return candidates[rand.Intn(len(candidates))], true
}

func (m *CarnivalMatch) genPosOccupied(mapInstance *Map, pos types.Point[int16], teamID constant.CarnivalTeam) bool {
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

func (m *CarnivalMatch) giveDebuff(resources *wz.Resources, target *Character, skillID uint32, skillLevel uint8) bool {
	if target == nil || resources == nil || skillID == 0 {
		return false
	}
	mobSkill := resources.GetMobSkill(skillID, skillLevel)
	if mobSkill == nil {
		return false
	}
	for _, flag := range constant.AllDebuffFlags() {
		if flag.DiseaseSkillID == uint16(skillID) {
			duration := time.Duration(mobSkill.DurationMs) * time.Millisecond
			if duration <= 0 {
				duration = 5 * time.Second
			}
			target.GiveDebuff(flag, duration, int16(mobSkill.X), uint16(skillID), uint16(skillLevel))
			return true
		}
	}
	return false
}
