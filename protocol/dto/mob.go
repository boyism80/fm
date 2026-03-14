package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MobStatusEntry struct {
	X       int16
	SkillID uint32
}

type Mob struct {
	OID        uint32
	MobId      uint32
	Position   types.Vector2[int16]
	Stance     uint8
	Foothold   int16
	Hp         uint16
	MaxHp      uint16
	Mp         uint16
	MaxMp      uint16
	StatusMask int32
	Statuses   []MobStatusEntry
}

func (m *Mob) Serialize(writer *stream.StreamWriter) error {
	if m == nil || m.StatusMask == 0 {
		writer.WriteU32(0)
		return nil
	}
	writer.Write32(m.StatusMask)
	for _, s := range m.Statuses {
		writer.Write16(s.X)
		if s.SkillID > 0 {
			writer.WriteU32(s.SkillID)
		}
		writer.Write16(32767)
	}
	return nil
}
