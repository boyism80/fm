package request

import (
	"github.com/boyism80/fm/stream"
)

type PetChat struct {
	Type   uint8
	Action uint8
	Text   string
}

func (*PetChat) Opcode() byte { return 0x85 }

func (p *PetChat) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.Type)
	writer.WriteU8(p.Action)
	writer.WriteStr16(p.Text)
	return nil
}

func (p *PetChat) Deserialize(reader *stream.StreamReader) {
	p.Type = reader.ReadU8()
	p.Action = reader.ReadU8()
	p.Text = reader.ReadStr16()
}
