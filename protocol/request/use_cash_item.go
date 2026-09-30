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

	if p.ItemID/10000 != 507 {
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
