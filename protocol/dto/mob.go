package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MobBuffEntry struct {
	X          int16
	SkillID    uint16
	SkillLevel uint16
}

type Mob struct {
	OID        uint32
	MobId      uint32
	Position   types.Vector2[int16]
	Stance     uint8
	Foothold   int16
	Hp         uint32
	MaxHp      uint32
	Mp         uint32
	MaxMp      uint32
	StatusMask int32
	Statuses   []MobBuffEntry
}

func (m *Mob) Serialize(writer *stream.StreamWriter) error {
	if m == nil || m.StatusMask == 0 {
		writer.WriteU32(0)
		return nil
	}
	writer.Write32(m.StatusMask)
	for _, s := range m.Statuses {
		writer.Write16(s.X)
		writer.WriteU16(s.SkillID)
		writer.WriteU16(s.SkillLevel)
		writer.Write16(32767)
	}
	return nil
}
