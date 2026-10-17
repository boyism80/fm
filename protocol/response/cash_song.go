package response

import (
	"github.com/boyism80/fm/stream"
)

type CashSong struct {
	ItemID uint32
	Name   string
}

func (p *CashSong) Opcode() uint16 {
	return 0x61
}

func (p *CashSong) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.ItemID)
	writer.WriteStr16(p.Name)
	return nil
}

func (p *CashSong) Deserialize(reader *stream.StreamReader) {
	p.ItemID = reader.ReadU32()
	p.Name = reader.ReadStr16()
}
