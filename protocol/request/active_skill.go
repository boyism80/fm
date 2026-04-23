package request

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

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
	MagnetMobData *MagnetMobData
	Position      *types.Vector2[int16]
}

func (*ActiveSkill) Opcode() byte { return 0x4A }

func (s *ActiveSkill) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (s *ActiveSkill) Deserialize(reader *stream.StreamReader) {
	s.OldX = reader.Read16()
	s.OldY = reader.Read16()
	s.SkillID = reader.ReadU32()
	s.SkillLevel = reader.ReadU8()

	available := reader.Remaining()

	switch s.SkillID {
	case 1121001, 1221001, 1321001:
		count := reader.ReadU32()
		s.MagnetMobData = &MagnetMobData{
			Mobs: make([]MagnetMobEntry, 0, count),
		}
		for i := uint32(0); i < count; i++ {
			var entry MagnetMobEntry
			entry.OID = reader.ReadU32()
			raw := reader.ReadU8()
			entry.Success = raw != 0
			s.MagnetMobData.Mobs = append(s.MagnetMobData.Mobs, entry)
		}
		s.MagnetMobData.Direction = reader.ReadU8()
	default:
		if available == 5 || available == 7 {
			s.Position = &types.Vector2[int16]{X: reader.Read16(), Y: reader.Read16()}
		}
	}
}
