package request

import "github.com/boyism80/fm/stream"

type UseCashItem struct {
	Slot   uint16
	ItemID uint32
	Text   string
	Ear    bool
	PetSN  uint64
	Target TeleportStoneTarget
}

func (*UseCashItem) Opcode() byte { return 0x3E }

func (p *UseCashItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(p.Slot)
	writer.WriteU32(p.ItemID)

	switch p.ItemID / 10000 {
	case 517:
		writer.WriteStr16(p.Text)
		return nil
	case 519:
		writer.WriteU64(p.PetSN)
		return nil
	case 504:
		p.Target.Serialize(writer)
		return nil
	case 507:
	default:
		return nil
	}
	switch p.ItemID % 10000 / 1000 {
	case 1, 8:
		writer.WriteStr16(p.Text)
	case 2:
		writer.WriteStr16(p.Text)
		writer.WriteBoolean(p.Ear)
	}
	return nil
}

func (p *UseCashItem) Deserialize(reader *stream.StreamReader) {
	p.Slot = reader.ReadU16()
	p.ItemID = reader.ReadU32()

	switch p.ItemID / 10000 {
	case 517:
		p.Text = reader.ReadStr16()
		return
	case 519:
		p.PetSN = reader.ReadU64()
		return
	case 504:
		p.Target.Deserialize(reader)
		return
	case 507:
	default:
		return
	}
	switch p.ItemID % 10000 / 1000 {
	case 1, 8:
		p.Text = reader.ReadStr16()
	case 2:
		p.Text = reader.ReadStr16()
		p.Ear = reader.ReadU8() != 0
	}
}
