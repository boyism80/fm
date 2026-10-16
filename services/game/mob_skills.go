package game

import (
	"fmt"
	"github.com/boyism80/fm/core/clock"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/wz"
	lua "github.com/yuin/gopher-lua"
)

type MobSkill struct {
	slot       wz.MobSkillSlot
	levelData  *wz.MobSkillLevelData
	lastUsedAt time.Time
}

func (s *MobSkill) IsOnCooldown(now time.Time) bool {
	if s == nil || s.levelData == nil {
		return true
	}
	if s.lastUsedAt.IsZero() {
		return false
	}
	if s.levelData.SummonOnce {
		return true
	}
	if s.levelData.CooltimeMs <= 0 {
		return false
	}
	cooldown := time.Duration(s.levelData.CooltimeMs) * time.Millisecond
	return now.Sub(s.lastUsedAt) < cooldown
}

const mobSkillCastWindow = 5 * time.Second

type MobSkills struct {
	owner     *Mob
	ordered   []*MobSkill
	byID      map[uint32]*MobSkill
	pending   *MobSkill
	pendingAt time.Time
}

func NewMobSkills(owner *Mob) *MobSkills {
	if owner == nil {
		panic("MobSkills owner is nil")
	}
	if owner.wz == nil {
		panic("MobSkills: mob model is nil")
	}
	if owner.gameWorld == nil {
		panic("MobSkills: game world is nil")
	}
	resources := owner.gameWorld.GetResources()
	if resources == nil {
		panic("MobSkills: resources are nil")
	}
	skills := &MobSkills{
		owner:   owner,
		ordered: make([]*MobSkill, 0, len(owner.wz.Skills)),
		byID:    make(map[uint32]*MobSkill, len(owner.wz.Skills)),
	}
	for _, slot := range owner.wz.Skills {
		if !owner.wz.HasSkill(slot.SkillID, slot.Level) {
			continue
		}
		levelData := resources.GetMobSkill(slot.SkillID, slot.Level)
		if levelData == nil {
			continue
		}
		entity := &MobSkill{
			slot:      slot,
			levelData: levelData,
		}
		skills.ordered = append(skills.ordered, entity)
		skills.byID[slot.SkillID] = entity
	}
	return skills
}

func (sc *MobSkills) Ordered() []*MobSkill {
	if sc == nil {
		return nil
	}
	return sc.ordered
}

func (sc *MobSkills) Get(skillID uint32, level uint8) *MobSkill {
	if sc == nil {
		return nil
	}
	entity := sc.byID[skillID]
	if entity == nil || entity.slot.Level != level {
		return nil
	}
	return entity
}

func (sc *MobSkills) MarkUsed(skillID uint32, now time.Time) {
	if sc == nil {
		return
	}
	entity := sc.byID[skillID]
	if entity == nil {
		return
	}
	entity.lastUsedAt = now
	if skillID == 140 || skillID == 141 {
		other := uint32(141)
		if skillID == 141 {
			other = 140
		}
		if sibling := sc.byID[other]; sibling != nil {
			sibling.lastUsedAt = now
		}
	}
}

func (sc *MobSkills) Choice(controller *Character) *MobSkill {
	if sc == nil {
		return nil
	}
	now := clock.Now()
	for _, skill := range sc.ordered {
		if skill == nil || skill.levelData == nil {
			continue
		}
		if skill.IsOnCooldown(now) {
			continue
		}
		maxHP := sc.owner.GetMaxHp()
		if maxHP == 0 {
			continue
		}
		currentPercent := int((uint64(sc.owner.GetHp()) * 100) / uint64(maxHP))
		if currentPercent > skill.levelData.HpPercent {
			continue
		}
		if skill.levelData.Limit > 0 {
			mapInstance := sc.owner.GetMap()
			if mapInstance == nil {
				continue
			}
			if len(mapInstance.GetMobs()) >= int(skill.levelData.Limit) {
				continue
			}
		}
		if !sc.chooseByScript(controller, skill) {
			continue
		}
		sc.MarkUsed(skill.slot.SkillID, now)
		sc.pending = skill
		sc.pendingAt = now
		return skill
	}
	return nil
}

// TakePending returns the skill Choice gave the controller to cast, once. The client names the skill in its cast
// packet, so anything else (another skill, a second cast, a cast long after the choice) is ignored.
func (sc *MobSkills) TakePending(skillID uint32, level uint8, now time.Time) *MobSkill {
	if sc == nil {
		return nil
	}
	skill := sc.pending
	if skill == nil || skill.slot.SkillID != skillID || skill.slot.Level != level || now.Sub(sc.pendingAt) > mobSkillCastWindow {
		return nil
	}
	sc.pending = nil

	if skill.levelData.Limit > 0 {
		mapInstance := sc.owner.GetMap()
		if mapInstance == nil || len(mapInstance.GetMobs()) >= int(skill.levelData.Limit) {
			return nil
		}
	}
	return skill
}

func (sc *MobSkills) chooseByScript(controller *Character, skill *MobSkill) bool {
	if sc == nil || sc.owner == nil || skill == nil || skill.levelData == nil {
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
	scriptPath := fmt.Sprintf("script/mob/skill/%d.lua", skill.slot.SkillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return true
	}
	ret, err := luax.Call(thread, "on_mob_skill_choose", sc.owner, skill)
	if err != nil || ret == nil || ret == lua.LNil {
		return true
	}
	return lua.LVAsBool(ret)
}

func (m *Mob) ConsumeSkillMP(skill *MobSkill) {
	if skill.levelData == nil || skill.levelData.MpCon <= 0 {
		return
	}
	cost := uint32(skill.levelData.MpCon)
	next := uint32(0)
	if current := m.GetMp(); current > cost {
		next = current - cost
	}
	m.SetMp(max(next, min(m.GetMaxMp(), 500)), false)
}
