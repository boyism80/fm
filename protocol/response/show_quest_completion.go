package response

import "github.com/boyism80/fm/stream"

type ShowQuestCompletion struct {
	QuestID uint16
}

func (p *ShowQuestCompletion) Opcode() uint16 {
	return 0x25
}

func (p *ShowQuestCompletion) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(p.QuestID)
	return nil
}

func (p *ShowQuestCompletion) Deserialize(reader *stream.StreamReader) {
}
