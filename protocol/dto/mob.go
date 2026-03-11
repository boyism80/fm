package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

// Mob represents mob data for protocol
type Mob struct {
	OID      uint32
	MobId    uint32
	Position types.Vector2[int16]
	Stance   uint8
	Foothold int16
	Hp       uint16
	MaxHp    uint16
	Mp       uint16
	MaxMp    uint16
}

// Serialize serializes mob data
func (m *Mob) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(0) // Control status
	return nil
}
