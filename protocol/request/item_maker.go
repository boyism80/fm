package request

import "github.com/boyism80/fm/stream"

type ItemMaker struct {
	Type      uint32
	ItemID    uint32
	Catalyst  bool
	Enhancers []uint32
	Reserved  uint32
	Slot      int32
	Malformed bool
}

func (*ItemMaker) Opcode() byte { return 0x60 }

func (p *ItemMaker) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Type)
	writer.WriteU32(p.ItemID)
	switch p.Type {
	case 1, 2:
		writer.WriteBoolean(p.Catalyst)
		writer.WriteU32(uint32(len(p.Enhancers)))
		for _, itemID := range p.Enhancers {
			writer.WriteU32(itemID)
		}
	case 4:
		writer.WriteU32(p.Reserved)
		writer.Write32(p.Slot)
	}
	return nil
}

func (p *ItemMaker) Deserialize(reader *stream.StreamReader) {
	p.Type = reader.ReadU32()
	p.ItemID = reader.ReadU32()
	switch p.Type {
	case 1, 2:
		p.Catalyst = reader.ReadU8() != 0
		count := reader.ReadU32()
		if count > 16 {
			p.Malformed = true
			return
		}
		p.Enhancers = make([]uint32, count)
		for i := range p.Enhancers {
			p.Enhancers[i] = reader.ReadU32()
		}
	case 3:
	case 4:
		p.Reserved = reader.ReadU32()
		p.Slot = int32(reader.ReadU32())
	default:
		p.Malformed = true
	}
}
