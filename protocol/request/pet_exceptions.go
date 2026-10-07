package request

import (
	"github.com/boyism80/fm/stream"
)

type PetExceptions struct {
	ItemIDs []uint32
}

func (*PetExceptions) Opcode() byte { return 0x89 }

func (p *PetExceptions) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(len(p.ItemIDs)))
	for _, id := range p.ItemIDs {
		writer.WriteU32(id)
	}
	return nil
}

func (p *PetExceptions) Deserialize(reader *stream.StreamReader) {
	count := reader.ReadU8()
	p.ItemIDs = make([]uint32, 0, count)
	for range count {
		p.ItemIDs = append(p.ItemIDs, reader.ReadU32())
	}
}
