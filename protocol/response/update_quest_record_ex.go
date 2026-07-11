package response

import "github.com/boyism80/fm/stream"

const QuestWireRecordExSubOpcode uint8 = 0x0A

type UpdateQuestRecordEx struct {
	QuestID uint16
	Data    string
}

func (p *UpdateQuestRecordEx) Opcode() uint16 {
	return 0x1C
}

func (p *UpdateQuestRecordEx) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(QuestWireRecordExSubOpcode)
	writer.WriteU16(p.QuestID)
	if p.Data == "" {
		writer.WriteStr8("")
	} else {
		writer.WriteStr8(p.Data)
	}
	return nil
}

func (p *UpdateQuestRecordEx) Deserialize(reader *stream.StreamReader) {
}
