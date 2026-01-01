package response

import (
	"github.com/boyism80/fm/stream"
)

// ShowForeignEffect represents the SHOW_FOREIGN_EFFECT packet (0x8F)
// This packet is used to show effects for other players on the map.
// EffectID 0 = Level up effect
type ShowForeignEffect struct {
	CharacterID uint32
	EffectID    uint8
}

func (p *ShowForeignEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowForeignEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(p.EffectID)
	return nil
}

func (p *ShowForeignEffect) Deserialize(reader *stream.StreamReader) error {
	return nil
}

