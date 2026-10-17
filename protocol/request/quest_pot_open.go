package request

import "github.com/boyism80/fm/stream"

type QuestPotOpen struct {
	QuestID uint16
}

func (*QuestPotOpen) Opcode() byte { return 0xC3 }

func (p *QuestPotOpen) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(p.QuestID)
	return nil
}

func (p *QuestPotOpen) Deserialize(reader *stream.StreamReader) {
	p.QuestID = reader.ReadU16()
}
