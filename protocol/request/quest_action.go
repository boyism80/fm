package request

import "github.com/boyism80/fm/stream"

type QuestActionMode uint8

const (
	QuestActionRestoreLostItem QuestActionMode = 0
	QuestActionStart           QuestActionMode = 1
	QuestActionComplete        QuestActionMode = 2
	QuestActionForfeit         QuestActionMode = 3
	QuestActionScriptedStart   QuestActionMode = 4
	QuestActionScriptedEnd     QuestActionMode = 5
)

type QuestAction struct {
	Mode      QuestActionMode
	QuestID   uint16
	NPCID     uint32
	ItemID    uint32
	Selection *uint32
}

func (*QuestAction) Opcode() byte { return 0x5A }

func (p *QuestAction) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	writer.WriteU16(p.QuestID)

	switch p.Mode {
	case QuestActionRestoreLostItem:
		writer.WriteU32(0)
		writer.WriteU32(p.ItemID)

	case QuestActionStart, QuestActionScriptedStart, QuestActionScriptedEnd:
		writer.WriteU32(p.NPCID)

	case QuestActionComplete:
		writer.WriteU32(p.NPCID)
		writer.WriteU32(0)
		selection := uint32(0)
		if p.Selection != nil {
			selection = *p.Selection
		}
		writer.WriteU32(selection)
	}
	return nil
}

func (p *QuestAction) Deserialize(reader *stream.StreamReader) {
	p.Mode = QuestActionMode(reader.ReadU8())
	p.QuestID = reader.ReadU16()

	switch p.Mode {
	case QuestActionRestoreLostItem:
		reader.ReadU32()
		p.ItemID = reader.ReadU32()

	case QuestActionStart, QuestActionScriptedStart, QuestActionScriptedEnd:
		p.NPCID = reader.ReadU32()

	case QuestActionComplete:
		p.NPCID = reader.ReadU32()
		reader.Skip(4)
		selection := reader.ReadU32()
		p.Selection = &selection

	case QuestActionForfeit:
	}
}
