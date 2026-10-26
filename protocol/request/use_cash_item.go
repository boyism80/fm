package request

import "github.com/boyism80/fm/stream"

type UseCashItem struct {
	Slot         uint16
	ItemID       uint32
	Text         string
	Ear          bool
	PetSN        uint64
	Target       TeleportStoneTarget
	SearchID     uint32
	HighFirst    bool
	Recipient    string
	Lines        []string
	Attached     bool
	AttachedType uint32
	AttachedSlot uint32
	// From and To are the stat flags of an AP reset, or the skills of an SP reset.
	From uint32
	To   uint32
	// TargetType and TargetSlot name the item a karma scissors, sealing lock or incubator works on;
	// EquipSlot is the worn equipment an item tag or magic hourglass works on.
	TargetType uint32
	TargetSlot uint32
	EquipSlot  int16
	Tick       uint32
	Malformed  bool
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
	case 523:
		writer.WriteU32(p.SearchID)
		writer.WriteBoolean(p.HighFirst)
		writer.WriteU32(0)
		return nil
	case 505:
		writer.WriteU32(p.To)
		writer.WriteU32(p.From)
		return nil
	case 506:
		if p.ItemID == 5060000 {
			writer.Write16(p.EquipSlot)
			return nil
		}
		writer.WriteU32(p.TargetType)
		writer.WriteU32(p.TargetSlot)
		return nil
	case 550:
		writer.Write16(p.EquipSlot)
		return nil
	case 552:
		writer.WriteU32(p.TargetType)
		writer.WriteU32(p.TargetSlot)
		return nil
	case 512:
		writer.WriteStr16(p.Text)
		return nil
	case 507:
		switch p.ItemID % 10000 / 1000 {
		case 0, 1:
			writer.WriteStr16(p.Text)
		case 2, 3, 4:
			writer.WriteStr16(p.Text)
			writer.WriteBoolean(p.Ear)
		case 6:
			writer.WriteStr16(p.Text)
			writer.WriteBoolean(p.Ear)
			writer.WriteBoolean(p.Attached)
			if p.Attached {
				writer.WriteU32(p.AttachedType)
				writer.WriteU32(p.AttachedSlot)
			}
		case 7:
			writer.WriteU8(uint8(len(p.Lines)))
			for _, line := range p.Lines {
				writer.WriteStr16(line)
			}
			writer.WriteBoolean(p.Ear)
		default:
			return nil
		}
	case 508, 537:
		writer.WriteStr16(p.Text)
	case 509:
		writer.WriteStr16(p.Recipient)
		writer.WriteStr16(p.Text)
	default:
		return nil
	}
	writer.WriteU32(p.Tick)
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
	case 523:
		p.SearchID = reader.ReadU32()
		p.HighFirst = reader.ReadU8() != 0
		return
	case 505:
		p.To = reader.ReadU32()
		p.From = reader.ReadU32()
		return
	case 506:
		if p.ItemID == 5060000 {
			p.EquipSlot = reader.Read16()
			return
		}
		p.TargetType = reader.ReadU32()
		p.TargetSlot = reader.ReadU32()
		return
	case 550:
		p.EquipSlot = reader.Read16()
		return
	case 552:
		p.TargetType = reader.ReadU32()
		p.TargetSlot = reader.ReadU32()
		return
	case 512:
		p.Text = reader.ReadStr16()
		return
	case 507:
		switch p.ItemID % 10000 / 1000 {
		case 0, 1:
			p.Text = reader.ReadStr16()
		case 2, 3, 4:
			p.Text = reader.ReadStr16()
			p.Ear = reader.ReadU8() != 0
		case 6:
			p.Text = reader.ReadStr16()
			p.Ear = reader.ReadU8() != 0
			p.Attached = reader.ReadU8() != 0
			if p.Attached {
				p.AttachedType = reader.ReadU32()
				p.AttachedSlot = reader.ReadU32()
			}
		case 7:
			count := reader.ReadU8()
			if count < 1 || count > 3 {
				p.Malformed = true
				return
			}
			p.Lines = make([]string, count)
			for i := range p.Lines {
				p.Lines[i] = reader.ReadStr16()
			}
			p.Ear = reader.ReadU8() != 0
		default:
			return
		}
	case 508, 537:
		p.Text = reader.ReadStr16()
	case 509:
		p.Recipient = reader.ReadStr16()
		p.Text = reader.ReadStr16()
	default:
		return
	}
	p.Tick = reader.ReadU32()
}
