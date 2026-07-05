package response

import "github.com/boyism80/fm/stream"

type UpdateQuestInfo struct {
	Progress    uint8
	QuestID     uint16
	NPCID       uint32
	NextQuestID uint32
}

func (p *UpdateQuestInfo) Opcode() uint16 {
	return 0x9C
}

func (p *UpdateQuestInfo) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.Progress)
	writer.WriteU16(p.QuestID)
	writer.WriteU32(p.NPCID)
	writer.WriteU32(p.NextQuestID)
	return nil
}

func (p *UpdateQuestInfo) Deserialize(reader *stream.StreamReader) {
}
