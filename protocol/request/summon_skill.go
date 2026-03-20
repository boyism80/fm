package request

import (
	"errors"
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

func (m *SummonSkill) Deserialize(reader *stream.StreamReader) error {
	var err error
	m.SummonOID, err = reader.ReadU32()
	if err != nil {
		return err
	}
	if reader.Remaining() < 4 {
		return nil
	}
	m.SubSkillID, err = reader.ReadU32()
	if err != nil {
		return err
	}
	switch m.SubSkillID {
	case skillIDBeholderBuff:
		if reader.Remaining() < 2 {
			return errors.New("summon_skill: beholder buff expects 2 bytes after sub_skill_id")
		}
		reader.Skip(1)
		m.BuffEffectIndex, err = reader.ReadU8()
		if err != nil {
			return err
		}
		if reader.Remaining() != 0 {
			return fmt.Errorf("summon_skill: %d trailing bytes after beholder buff payload", reader.Remaining())
		}
	case skillIDBeholderHealing:
		if reader.Remaining() == 1 {
			reader.Skip(1)
		}
		if reader.Remaining() != 0 {
			return fmt.Errorf("summon_skill: %d trailing bytes after beholder healing sub_skill_id", reader.Remaining())
		}
	default:
		if reader.Remaining() != 0 {
			return fmt.Errorf("summon_skill: %d trailing bytes after unknown sub_skill_id", reader.Remaining())
		}
	}
	return nil
}
