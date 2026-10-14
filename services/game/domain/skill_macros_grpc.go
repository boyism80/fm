package domain

import internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"

func NewSkillMacrosFromInternalProto(entries []*internal.SkillMacro) []SkillMacro {
	macros := make([]SkillMacro, 0, min(len(entries), 5))
	for _, entry := range entries {
		if entry == nil || entry.GetPosition() >= 5 {
			continue
		}
		macro := SkillMacro{
			Name:  entry.GetName(),
			Shout: entry.GetShout(),
			Skills: [3]uint32{
				entry.GetSkill1(),
				entry.GetSkill2(),
				entry.GetSkill3(),
			},
		}
		for len(macros) < int(entry.GetPosition()) {
			macros = append(macros, SkillMacro{})
		}
		if int(entry.GetPosition()) == len(macros) {
			macros = append(macros, macro)
		} else {
			macros[entry.GetPosition()] = macro
		}
	}
	return macros
}

func (m *SkillMacros) BuildRPC() []*internal.SkillMacro {
	entries := make([]*internal.SkillMacro, 0, len(m.entries))
	for position, entry := range m.entries {
		entries = append(entries, &internal.SkillMacro{
			Position: uint32(position),
			Name:     entry.Name,
			Shout:    entry.Shout,
			Skill1:   entry.Skills[0],
			Skill2:   entry.Skills[1],
			Skill3:   entry.Skills[2],
		})
	}
	return entries
}
