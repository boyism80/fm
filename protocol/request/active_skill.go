package request

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type ActiveSkill struct {
	OldX       int16
	OldY       int16
	SkillID    uint32
	SkillLevel uint8

	// Skill-specific data (varies by SkillID)
	// For Rush skills (1121001, 1221001, 1321001):
	MobID     uint32
	Success   byte
	Direction byte

	// For other skills (optional position):
	Position *types.Vector2[int16]
}

func (s *ActiveSkill) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (s *ActiveSkill) Deserialize(reader *stream.StreamReader) error {
	var err error
	if s.OldX, err = reader.Read16(); err != nil {
		return err
	}
	if s.OldY, err = reader.Read16(); err != nil {
		return err
	}
	if s.SkillID, err = reader.ReadU32(); err != nil {
		return err
	}
	if s.SkillLevel, err = reader.ReadU8(); err != nil {
		return err
	}

	available := reader.Remaining()

	// Rush skills (1121001, 1221001, 1321001)
	if s.SkillID == 1121001 || s.SkillID == 1221001 || s.SkillID == 1321001 {
		if available >= 6 {
			if s.MobID, err = reader.ReadU32(); err != nil {
				return err
			}
			if s.Success, err = reader.ReadU8(); err != nil {
				return err
			}
			if s.Direction, err = reader.ReadU8(); err != nil {
				return err
			}
		}
	} else if s.SkillID == 9001004 {
		// Hidden toggle - no additional data
	} else {
		// Other skills - optional position
		if available == 5 || available == 7 {
			var x, y int16
			if x, err = reader.Read16(); err != nil {
				return err
			}
			if y, err = reader.Read16(); err != nil {
				return err
			}
			s.Position = &types.Vector2[int16]{X: x, Y: y}
		}
	}

	return nil
}
