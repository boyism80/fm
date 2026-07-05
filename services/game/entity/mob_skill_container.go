package entity

import (
	"fmt"
	"github.com/boyism80/fm/core/clock"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

type MobSkill struct {
	Slot       wz.MobSkillSlot
	LevelData  *wz.MobSkillLevelData
	LastUsedAt time.Time
}

func (s *MobSkill) IsOnCooldown(now time.Time) bool {
	if s == nil || s.LevelData == nil {
		return true
	}
	if s.LastUsedAt.IsZero() {
		return false
	}
	if s.LevelData.SummonOnce {
		return true
	}
	if s.LevelData.CooltimeMs <= 0 {
		return false
	}
	cooldown := time.Duration(s.LevelData.CooltimeMs) * time.Millisecond
	return now.Sub(s.LastUsedAt) < cooldown
}

type MobSkillContainer struct {
	owner   *Mob
	ordered []*MobSkill
	byID    map[uint32]*MobSkill
}

func NewMobSkillContainer(owner *Mob) *MobSkillContainer {
	if owner == nil {
		panic("MobSkillContainer owner is nil")
	}
	if owner.Wz == nil {
		panic("MobSkillContainer: mob model is nil")
	}
	if owner.GameWorld == nil {
		panic("MobSkillContainer: game world is nil")
	}
	resources := owner.GameWorld.GetResources()
	if resources == nil {
		panic("MobSkillContainer: resources are nil")
	}
	container := &MobSkillContainer{
		owner:   owner,
		ordered: make([]*MobSkill, 0, len(owner.Wz.Skills)),
		byID:    make(map[uint32]*MobSkill, len(owner.Wz.Skills)),
	}
	for _, slot := range owner.Wz.Skills {
		if !owner.Wz.HasSkill(slot.SkillID, slot.Level) {
			continue
		}
		levelData := resources.GetMobSkill(slot.SkillID, slot.Level)
		if levelData == nil {
			continue
		}
		entity := &MobSkill{
			Slot:      slot,
			LevelData: levelData,
		}
		container.ordered = append(container.ordered, entity)
		container.byID[slot.SkillID] = entity
	}
	return container
}

func (sc *MobSkillContainer) Ordered() []*MobSkill {
	if sc == nil {
		return nil
	}
	return sc.ordered
}

func (sc *MobSkillContainer) Get(skillID uint32, level uint8) *MobSkill {
	if sc == nil {
		return nil
	}
	entity := sc.byID[skillID]
	if entity == nil || entity.Slot.Level != level {
		return nil
	}
	return entity
}

func (sc *MobSkillContainer) MarkUsed(skillID uint32, now time.Time) {
	if sc == nil {
		return
	}
	entity := sc.byID[skillID]
	if entity == nil {
		return
	}
	entity.LastUsedAt = now
	if skillID == 140 || skillID == 141 {
		other := uint32(141)
		if skillID == 141 {
			other = 140
		}
		if sibling := sc.byID[other]; sibling != nil {
			sibling.LastUsedAt = now
		}
	}
}

func (sc *MobSkillContainer) Choice(controller *Character) *MobSkill {
	if sc == nil {
		return nil
	}
	now := clock.Now()
	for _, skill := range sc.ordered {
		if skill == nil || skill.LevelData == nil {
			continue
		}
		if skill.IsOnCooldown(now) {
			continue
		}
		if !sc.meetsCommonConditions(skill) {
			continue
		}
		if !sc.chooseByScript(controller, skill) {
			continue
		}
		sc.MarkUsed(skill.Slot.SkillID, now)
		return skill
	}
	return nil
}

func (sc *MobSkillContainer) meetsCommonConditions(skill *MobSkill) bool {
	if sc == nil || sc.owner == nil || skill == nil || skill.LevelData == nil {
		return false
	}
	maxHP := sc.owner.GetMaxHp()
	if maxHP == 0 {
		return false
	}
	currentPercent := int((uint64(sc.owner.GetHp()) * 100) / uint64(maxHP))
	if currentPercent > skill.LevelData.HpPercent {
		return false
	}
	if skill.LevelData.Limit > 0 {
		mapInstance := sc.owner.GetMap()
		if mapInstance == nil {
			return false
		}
		if len(mapInstance.GetMobs()) >= int(skill.LevelData.Limit) {
			return false
		}
	}
	return true
}

func (sc *MobSkillContainer) chooseByScript(controller *Character, skill *MobSkill) bool {
	if sc == nil || sc.owner == nil || skill == nil || skill.LevelData == nil {
		return false
	}
	mapInstance := sc.owner.GetMap()
	if mapInstance == nil {
		return false
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return true
	}
	scriptPath := fmt.Sprintf("script/mob/skill/%d.lua", skill.Slot.SkillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return true
	}
	hook := fmt.Sprintf("on_mob_skill_choose_%d", skill.Slot.SkillID)
	chosen := true
	luax.CallAsync(root, thread, hook, sc.owner, skill).Then(func(value interface{}) (interface{}, error) {
		vals := luax.ResultValues(value)
		if len(vals) == 0 || vals[0] == nil {
			chosen = true
			return nil, nil
		}
		chosen = lua.LVAsBool(vals[0])
		return nil, nil
	}).OnError(func(err error) {
		chosen = true
	})
	return chosen
}
