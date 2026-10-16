package game

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (sc *Skills) Load(skills []*internal.Skill) {
	for _, pb := range skills {
		if pb == nil {
			continue
		}
		entry, err := NewSkillEntryFromInternalProto(sc.owner, pb, sc.owner.gameWorld)
		if err != nil {
			continue
		}
		sc.Bind(pb.GetSkillId(), entry)
	}
}

func (sc *Skills) BuildRPC() []*internal.Skill {
	skills := make([]*internal.Skill, 0, 64)
	sc.ForEach(func(skillID uint32, entry *SkillEntry) {
		if entry == nil {
			return
		}
		if pb := entry.BuildRPC(sc.owner.GetID(), skillID); pb != nil {
			skills = append(skills, pb)
		}
	})
	return skills
}
