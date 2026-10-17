package request

import "github.com/boyism80/fm/stream"

type UseRewardBox struct {
	Slot   int16
	ItemID uint32
}

func (*UseRewardBox) Opcode() byte { return 0x5F }

func (p *UseRewardBox) Serialize(writer *stream.StreamWriter) error {
	writer.Write16(p.Slot)
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *UseRewardBox) Deserialize(reader *stream.StreamReader) {
	p.Slot = reader.Read16()
	p.ItemID = reader.ReadU32()
}
