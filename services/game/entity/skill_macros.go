package entity

import "unicode/utf8"

type SkillMacro struct {
	Name   string
	Shout  bool
	Skills [3]uint32
}

type SkillMacros struct {
	owner   *Character
	entries []SkillMacro
}

func NewSkillMacros(owner *Character) *SkillMacros {
	return &SkillMacros{owner: owner}
}

func (m *SkillMacros) Entries() []SkillMacro {
	return append([]SkillMacro(nil), m.entries...)
}

func (m *SkillMacros) Load(entries []SkillMacro) {
	m.entries = append([]SkillMacro(nil), entries...)
}

func (m *SkillMacros) Replace(entries []SkillMacro) bool {
	if len(entries) > 5 {
		return false
	}
	for _, entry := range entries {
		if utf8.RuneCountInString(entry.Name) > 30 {
			return false
		}
		for _, skillID := range entry.Skills {
			if skillID == 0 {
				continue
			}
			skill := m.owner.Skills.Learned(skillID)
			if skill == nil || skill.Level() < 1 {
				return false
			}
		}
	}
	m.entries = append([]SkillMacro(nil), entries...)
	return true
}
