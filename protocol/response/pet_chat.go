package response

import (
	"github.com/boyism80/fm/stream"
)

type PetChat struct {
	CharacterID uint32
	Type        uint8
	Action      uint8
	Text        string
}

func (p *PetChat) Opcode() uint16 {
	return 0x77
}

func (p *PetChat) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(p.Type)
	writer.WriteU8(p.Action)
	writer.WriteStr16(p.Text)
	writer.WriteU8(0)
	return nil
}

func (p *PetChat) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.Type = reader.ReadU8()
	p.Action = reader.ReadU8()
	p.Text = reader.ReadStr16()
}
