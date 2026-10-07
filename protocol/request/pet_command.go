package request

import (
	"github.com/boyism80/fm/stream"
)

type PetCommand struct {
	CalledByName bool
	Index        uint8
}

func (*PetCommand) Opcode() byte { return 0x86 }

func (p *PetCommand) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(p.CalledByName)
	writer.WriteU8(p.Index)
	return nil
}

func (p *PetCommand) Deserialize(reader *stream.StreamReader) {
	p.CalledByName = reader.ReadBool()
	p.Index = reader.ReadU8()
}
