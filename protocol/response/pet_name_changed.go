package response

import (
	"github.com/boyism80/fm/stream"
)

type PetNameChanged struct {
	CharacterID uint32
	Name        string
}

func (p *PetNameChanged) Opcode() uint16 {
	return 0x78
}

func (p *PetNameChanged) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteStr16(p.Name)
	return nil
}

func (p *PetNameChanged) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.Name = reader.ReadStr16()
}
