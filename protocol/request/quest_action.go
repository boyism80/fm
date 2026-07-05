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
		if reader.Remaining() >= 4 {
			reader.ReadU32()
		}
		if reader.Remaining() >= 4 {
			selection := reader.ReadU32()
			p.Selection = &selection
		} else if reader.Remaining() >= 2 {
			selection := uint32(reader.ReadU16())
			p.Selection = &selection
		}

	case QuestActionForfeit:
	}
}
