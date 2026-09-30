package request

import "github.com/boyism80/fm/stream"

type UseCashItem struct {
	Slot   uint16
	ItemID uint32
	Text   string
	Ear    bool
}

func (*UseCashItem) Opcode() byte { return 0x3E }

func (p *UseCashItem) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseCashItem) Deserialize(reader *stream.StreamReader) {
	p.Slot = reader.ReadU16()
	p.ItemID = reader.ReadU32()
	switch p.ItemID {
	case 5070000, 5071000:
		p.Text = reader.ReadStr16()
	case 5072000, 5073000, 5074000:
		p.Text = reader.ReadStr16()
		p.Ear = reader.ReadU8() != 0
	}
}
