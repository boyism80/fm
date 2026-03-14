package request

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

// MagnetMobEntry holds the result for a single mob pulled by Monster Magnet.
type MagnetMobEntry struct {
	OID     uint32
	Success bool
}

type MagnetMobData struct {
	Mobs      []MagnetMobEntry
	Direction byte
}

type ActiveSkill struct {
	OldX          int16
	OldY          int16
	SkillID       uint32
	SkillLevel    uint8
	MagnetMobData MagnetMobData
	Position      *types.Vector2[int16]
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

	switch s.SkillID {
	case 1121001, 1221001, 1321001:
		var count uint32
		if count, err = reader.ReadU32(); err != nil {
			return err
		}
		s.MagnetMobData.Mobs = make([]MagnetMobEntry, 0, count)
		for i := uint32(0); i < count; i++ {
			var entry MagnetMobEntry
			if entry.OID, err = reader.ReadU32(); err != nil {
				return err
			}
			var raw uint8
			if raw, err = reader.ReadU8(); err != nil {
				return err
			}
			entry.Success = raw != 0
			s.MagnetMobData.Mobs = append(s.MagnetMobData.Mobs, entry)
		}
		if s.MagnetMobData.Direction, err = reader.ReadU8(); err != nil {
			return err
		}
	default:
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
