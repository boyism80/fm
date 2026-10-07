package response

import (
	"github.com/boyism80/fm/stream"
)

type PetCommand struct {
	CharacterID uint32
	Food        bool
	Index       uint8
	Success     bool
}

func (p *PetCommand) Opcode() uint16 {
	return 0x7A
}

func (p *PetCommand) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteBoolean(p.Food)
	if p.Food {
		writer.WriteBoolean(p.Success)
		return nil
	}
	writer.WriteU8(p.Index)
	writer.WriteBoolean(p.Success)
	return nil
}

func (p *PetCommand) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.Food = reader.ReadBool()
	if p.Food {
		p.Success = reader.ReadBool()
		return
	}
	p.Index = reader.ReadU8()
	p.Success = reader.ReadBool()
}
