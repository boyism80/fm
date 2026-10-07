package response

import (
	"github.com/boyism80/fm/stream"
)

type PetExceptions struct {
	CharacterID uint32
	SN          uint64
	ItemIDs     []uint32
}

func (p *PetExceptions) Opcode() uint16 {
	return 0x79
}

func (p *PetExceptions) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU64(p.SN)
	writer.WriteU8(uint8(len(p.ItemIDs)))
	for _, id := range p.ItemIDs {
		writer.WriteU32(id)
	}
	return nil
}

func (p *PetExceptions) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.SN = reader.ReadU64()
	count := reader.ReadU8()
	p.ItemIDs = make([]uint32, 0, count)
	for range count {
		p.ItemIDs = append(p.ItemIDs, reader.ReadU32())
	}
}
