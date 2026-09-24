package entity

import (
	"fmt"
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
		spend := entry.SpendCP
		if spend <= 0 {
			spend = skillDef.SpendCP
		}
		if spend <= 0 {
			return false
		}
		enemy := m.enemyTeam(team.TeamID)
		if enemy == nil {
			return false
		}
		if !m.debuffEnemies(gw, resources, enemy, entry.ID, skillDef) {
			return false
		}
		if !team.UseCP(ch, spend) {
			return false
		}
		ch.Listener.OnCarnivalSummon(ch, tab, uint8(idx), ch.GetName())
		return true
	case 2:
		guardianDef := resources.GetMCGuardian(uint32(num))
		if guardianDef == nil {
			return false
		}
		if guardianDef.SpendCP <= 0 || !team.UseCP(ch, guardianDef.SpendCP) {
			return false
		}
		reactorName := fmt.Sprintf("%d%d", team.TeamID, num)
		if existing := mapInstance.ReactorByName(reactorName); existing != nil && existing.State < 5 {
			return false
		}
		pos, ok := pickCarnivalGenPos(mc.GuardianGenPos, team.TeamID, mapInstance, false)
		if !ok {
			return false
		}
		reactorID := 9980000 + uint32(team.TeamID)
		spawned, err := mapInstance.SpawnReactorAtNamed(reactorID, pos, reactorName)
		if err != nil || spawned == nil {
			return false
		}
		m.buffAllyMobs(mapInstance, team.TeamID, resources, guardianDef.MobSkillID, guardianDef.Level)
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

func (m *CarnivalMatch) debuffEnemies(gw GameWorld, resources *wz.Resources, enemy *CarnivalTeam, catalogID uint32, skillDef *wz.MCSkill) bool {
	if m == nil || enemy == nil || resources == nil || skillDef == nil {
		return false
	}
	targets := enemy.Members(gw)
	if len(targets) == 0 {
		return false
	}
	rand.Shuffle(len(targets), func(i, j int) {
		targets[i], targets[j] = targets[j], targets[i]
	})
	chance := GlobalCarnivalRegistry().SkillHitChance(catalogID, skillDef.HitChance)
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
	flag, ok := carnivalMobSkillBuffFlag(skillID)
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

func carnivalMobSkillBuffFlag(skillID uint32) (constant.MobBuffFlag, bool) {
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
		resources := ch.GameWorld.GetResources()
		if resources != nil {
			skillDef := resources.GetMCSkill(wzConsume.NuffSkillID)
			if skillDef != nil {
				enemy := match.enemyTeam(team.TeamID)
				if enemy != nil {
					match.debuffEnemies(ch.GameWorld, resources, enemy, wzConsume.NuffSkillID, skillDef)
				}
			}
		}
	}
	return true
}
