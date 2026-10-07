package response

import (
	"github.com/boyism80/fm/stream"
)

type PetAutoHP struct {
	ItemID uint32
}

func (p *PetAutoHP) Opcode() uint16 {
	return 0xFE
}

func (p *PetAutoHP) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *PetAutoHP) Deserialize(reader *stream.StreamReader) {
	p.ItemID = reader.ReadU32()
}

type PetAutoMP struct {
	ItemID uint32
}

func (p *PetAutoMP) Opcode() uint16 {
	return 0xFF
}

func (p *PetAutoMP) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *PetAutoMP) Deserialize(reader *stream.StreamReader) {
	p.ItemID = reader.ReadU32()
}
