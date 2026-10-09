package response

import (
	"github.com/boyism80/fm/stream"
)

type ItemExpired struct {
	ItemIDs []uint32
}

func (p *ItemExpired) Opcode() uint16 {
	return 0x1C
}

func (p *ItemExpired) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(8)
	writer.WriteU8(uint8(len(p.ItemIDs)))
	for _, id := range p.ItemIDs {
		writer.WriteU32(id)
	}
	return nil
}

func (p *ItemExpired) Deserialize(reader *stream.StreamReader) {
}

type CashItemExpired struct {
	ItemID uint32
}

func (p *CashItemExpired) Opcode() uint16 {
	return 0x1C
}

func (p *CashItemExpired) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(2)
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *CashItemExpired) Deserialize(reader *stream.StreamReader) {
}
