package request

import (
	"github.com/boyism80/fm/stream"
)

type UseReturnScroll struct {
	Tick   uint32
	Slot   uint16
	ItemID uint32
}

func (*UseReturnScroll) Opcode() byte { return 0x44 }

func (p *UseReturnScroll) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseReturnScroll) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Slot = reader.ReadU16()
	p.ItemID = reader.ReadU32()
}
