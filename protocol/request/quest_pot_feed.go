package request

import "github.com/boyism80/fm/stream"

type QuestPotFeed struct {
	Slot    int16
	ItemID  uint32
	QuestID uint32
	Exp     uint32
}

func (*QuestPotFeed) Opcode() byte { return 0xC5 }

func (p *QuestPotFeed) Serialize(writer *stream.StreamWriter) error {
	writer.Write16(p.Slot)
	writer.WriteU32(p.ItemID)
	writer.WriteU32(p.QuestID)
	writer.WriteU32(p.Exp)
	return nil
}

func (p *QuestPotFeed) Deserialize(reader *stream.StreamReader) {
	p.Slot = reader.Read16()
	p.ItemID = reader.ReadU32()
	p.QuestID = reader.ReadU32()
	p.Exp = reader.ReadU32()
}
