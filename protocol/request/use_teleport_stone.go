package request

import "github.com/boyism80/fm/stream"

type UseTeleportStone struct {
	Slot   uint16
	ItemID uint32
	Target TeleportStoneTarget
}

func (*UseTeleportStone) Opcode() byte { return 0x43 }

func (p *UseTeleportStone) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(p.Slot)
	writer.WriteU32(p.ItemID)
	p.Target.Serialize(writer)
	return nil
}

func (p *UseTeleportStone) Deserialize(reader *stream.StreamReader) {
	p.Slot = reader.ReadU16()
	p.ItemID = reader.ReadU32()
	p.Target.Deserialize(reader)
}
