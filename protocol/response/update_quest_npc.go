package response

import "github.com/boyism80/fm/stream"

type UpdateQuestNPC struct {
	Progress    uint8
	QuestID     uint16
	NPCID       uint32
	NextQuestID uint32
}

func (p *UpdateQuestNPC) Opcode() uint16 {
	return 0x9C
}

func (p *UpdateQuestNPC) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.Progress)
	writer.WriteU16(p.QuestID)
	writer.WriteU32(p.NPCID)
	writer.WriteU32(p.NextQuestID)
	return nil
}

func (p *UpdateQuestNPC) Deserialize(reader *stream.StreamReader) {
}
