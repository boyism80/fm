package request

import (
	"github.com/boyism80/fm/stream"
)

// UseItem represents the USE_ITEM packet (0x37)
type UseItem struct {
	Tick   uint32
	Slot   uint16
	ItemID uint32
}

func (p *UseItem) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseItem) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Slot = reader.ReadU16()
	p.ItemID = reader.ReadU32()

}
