package request

import (
	"fmt"

	"github.com/boyism80/fm/stream"
)

const (
	skillIDBeholderHealing uint32 = 1320008
	skillIDBeholderBuff    uint32 = 1320009
)

type SummonSkill struct {
	SummonOID       uint32
	SubSkillID      uint32
	BuffEffectIndex uint8
}

func (m *SummonSkill) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *SummonSkill) Deserialize(reader *stream.StreamReader) {
	m.SummonOID = reader.ReadU32()
	if reader.Remaining() < 4 {
	}
	m.SubSkillID = reader.ReadU32()
	switch m.SubSkillID {
	case skillIDBeholderBuff:
		if reader.Remaining() < 2 {
			panic("summon_skill: beholder buff expects 2 bytes after sub_skill_id")
		}
		reader.Skip(1)
		m.BuffEffectIndex = reader.ReadU8()
		if reader.Remaining() != 0 {
			panic(fmt.Errorf("summon_skill: %d trailing bytes after beholder buff payload", reader.Remaining()))
		}
	case skillIDBeholderHealing:
		if reader.Remaining() == 1 {
			reader.Skip(1)
		}
		if reader.Remaining() != 0 {
			panic(fmt.Errorf("summon_skill: %d trailing bytes after beholder healing sub_skill_id", reader.Remaining()))
		}
	default:
		if reader.Remaining() != 0 {
			panic(fmt.Errorf("summon_skill: %d trailing bytes after unknown sub_skill_id", reader.Remaining()))
		}
	}
}
