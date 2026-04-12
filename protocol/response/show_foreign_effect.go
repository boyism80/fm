package response

import (
	"github.com/boyism80/fm/stream"
)

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

func (p *ShowForeignEffect) Deserialize(reader *stream.StreamReader) {
}
