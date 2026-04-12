package entity

import (
	"time"

	"github.com/boyism80/fm/game/wz"
)

type SkillEntry struct {
	Wz          *wz.Skill
	level       int
	MasterLevel int
	Expiration  time.Time
	CooldownEnd *time.Time
	Owner       *Character
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
		Wz:          w,
		level:       level,
		MasterLevel: masterLevel,
		Owner:       owner,
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
		s.applyPassiveAfterLevelChange(prev, level)
	}
	if s.Owner != nil {
		s.Owner.Listener.OnUpdateSkill(s.Owner, s.Wz.ID, int32(s.level), int32(s.MasterLevel))
	}
	if level == 0 {
		delete(s.Owner.Skills.entries, s.Wz.ID)
	}
}

func (s *SkillEntry) SetLevelAndMaster(level, masterLevel int) {
	if level < 0 {
		level = 0
	}
	if masterLevel < 0 {
		masterLevel = 0
	}
	s.MasterLevel = masterLevel
	prev := s.level
	s.level = level
	if prev != level {
		s.applyPassiveAfterLevelChange(prev, level)
	}
	if s.Owner != nil {
		s.Owner.Listener.OnUpdateSkill(s.Owner, s.Wz.ID, int32(s.level), int32(s.MasterLevel))
	}
	if level == 0 {
		delete(s.Owner.Skills.entries, s.Wz.ID)
	}
}

func (s *SkillEntry) SetMasterLevel(masterLevel int) {
	if masterLevel < 0 {
		masterLevel = 0
	}
	s.MasterLevel = masterLevel
	if s.Owner != nil {
		s.Owner.Listener.OnUpdateSkill(s.Owner, s.Wz.ID, int32(s.level), int32(s.MasterLevel))
	}
}

func (s *SkillEntry) applyPassiveAfterLevelChange(prevLevel, newLevel int) {
	if prevLevel == newLevel {
		return
	}
	skillID := s.Wz.ID
	if s.Owner == nil {
		return
	}
	if prevLevel > 0 && newLevel == 0 {
		s.Owner.Listener.OnSkillPassiveHook(s.Owner, skillID, "on_unpassive")
		return
	}
	if newLevel > 0 {
		s.Owner.Listener.OnSkillPassiveHook(s.Owner, skillID, "on_passive")
	}
}

func (s *SkillEntry) IsCooling() bool {
	if s.CooldownEnd == nil {
		return false
	}
	return time.Now().Before(*s.CooldownEnd)
}

func (s *SkillEntry) StartCooldown(duration time.Duration) {
	end := time.Now().Add(duration)
	s.CooldownEnd = &end
	sec := min(int(duration.Seconds()), 65535)
	s.notifyCooldown(uint16(sec))
}

func (s *SkillEntry) CooldownRemaining() time.Duration {
	if s.CooldownEnd == nil {
		return 0
	}
	if !time.Now().Before(*s.CooldownEnd) {
		return 0
	}
	return time.Until(*s.CooldownEnd)
}

func (s *SkillEntry) ClearCooldown() {
	s.CooldownEnd = nil
	s.notifyCooldown(0)
}

func (s *SkillEntry) notifyCooldown(remainingSec uint16) {
	if s.Owner != nil {
		s.Owner.Listener.OnSkillCooldown(s.Owner, s.Wz.ID, remainingSec)
	}
}
