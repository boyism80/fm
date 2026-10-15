package request

import (
	"github.com/boyism80/fm/stream"
)

type LieDetectorItem struct {
	Target string
	Slot   int16
	ItemID uint32
}

func (*LieDetectorItem) Opcode() byte { return 0x56 }

func (p *LieDetectorItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(p.Target)
	writer.Write16(p.Slot)
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *LieDetectorItem) Deserialize(reader *stream.StreamReader) {
	p.Target = reader.ReadStr16()
	p.Slot = reader.Read16()
	p.ItemID = reader.ReadU32()
}
