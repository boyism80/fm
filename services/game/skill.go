package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/wz"
)

type SkillEntry struct {
	wz          *wz.Skill
	level       int
	masterLevel int
	expiration  time.Time
	cooldownEnd *time.Time
	owner       *Character
}

func NewSkillEntry(owner *Character, w *wz.Skill, level, masterLevel int) *SkillEntry {
	if w == nil {
		panic("NewSkillEntry: nil wz.Skill")
	}
	if owner == nil {
		panic("NewSkillEntry: nil Owner")
	}
	if level < 0 {
		level = 0
	}
	if masterLevel < 0 {
		masterLevel = 0
	}
	return &SkillEntry{
		wz:          w,
		level:       level,
		masterLevel: masterLevel,
		owner:       owner,
	}
}

func (s *SkillEntry) Level() int {
	return s.level
}

func (s *SkillEntry) SetLevel(level int) {
	if level < 0 {
		level = 0
	}
	prev := s.level
	s.level = level
	if prev != level {
		s.refreshPassiveStats(prev, level)
	}
	if s.owner != nil {
		s.owner.listener.OnUpdateSkill(s.owner, s.wz.ID, int32(s.level), int32(s.masterLevel))
	}
	if level == 0 && s.masterLevel == 0 {
		delete(s.owner.skills.entries, s.wz.ID)
	}
}

func (s *SkillEntry) SetLevels(level, masterLevel int) {
	if level < 0 {
		level = 0
	}
	if masterLevel < 0 {
		masterLevel = 0
	}
	s.masterLevel = masterLevel
	prev := s.level
	s.level = level
	if prev != level {
		s.refreshPassiveStats(prev, level)
	}
	if s.owner != nil {
		s.owner.listener.OnUpdateSkill(s.owner, s.wz.ID, int32(s.level), int32(s.masterLevel))
	}
	if level == 0 && masterLevel == 0 {
		delete(s.owner.skills.entries, s.wz.ID)
	}
}

func (s *SkillEntry) SetMasterLevel(masterLevel int) {
	if masterLevel < 0 {
		masterLevel = 0
	}
	s.masterLevel = masterLevel
	if s.owner != nil {
		s.owner.listener.OnUpdateSkill(s.owner, s.wz.ID, int32(s.level), int32(s.masterLevel))
	}
}

func (s *SkillEntry) refreshPassiveStats(prevLevel, newLevel int) {
	if prevLevel == newLevel {
		return
	}
	if s.owner == nil {
		return
	}
	if prevLevel > 0 && newLevel == 0 {
		s.owner.skills.CallHook(nil, s, "on_unpassive")
		return
	}
	if newLevel > 0 {
		s.owner.skills.CallHook(nil, s, "on_passive")
	}
}

func (s *SkillEntry) IsCooling() bool {
	if s.cooldownEnd == nil {
		return false
	}
	return clock.Now().Before(*s.cooldownEnd)
}

func (s *SkillEntry) StartCooldown(duration time.Duration) {
	end := clock.Now().Add(duration)
	s.cooldownEnd = &end
	sec := min(int(duration.Seconds()), 65535)
	if s.owner != nil {
		s.owner.listener.OnSkillCooldown(s.owner, s.wz.ID, uint16(sec))
	}
}

func (s *SkillEntry) CooldownRemaining() time.Duration {
	if s.cooldownEnd == nil {
		return 0
	}
	if !clock.Now().Before(*s.cooldownEnd) {
		return 0
	}
	return time.Until(*s.cooldownEnd)
}

func (s *SkillEntry) ClearCooldown() {
	s.cooldownEnd = nil
	if s.owner != nil {
		s.owner.listener.OnSkillCooldown(s.owner, s.wz.ID, 0)
	}
}
