package response

import (
	"time"

	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

// UpdateQuestPlayTime is the binary PTR record of a daily play time quest: today's seconds and the days counted.
type UpdateQuestPlayTime struct {
	QuestID uint16
	Today   uint16
	Days    uint8
	At      time.Time
}

func (p *UpdateQuestPlayTime) Opcode() uint16 {
	return 0x1C
}

func (p *UpdateQuestPlayTime) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(QuestWireRecordExSubOpcode)
	writer.WriteU16(p.QuestID)
	writer.WriteU16(0x15)
	writer.Write([]byte("PTR="))
	writer.WriteU8(1)
	writer.WriteU8(1)
	writer.WriteU64(util.ToFileTime(p.At))
	writer.WriteU8(1)
	writer.WriteU8(0x12)
	writer.WriteU8(0xCD)
	writer.WriteU16(p.Today)
	writer.WriteU8(p.Days)
	writer.WriteU8(0xCD)
	return nil
}

func (p *UpdateQuestPlayTime) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.QuestID = reader.ReadU16()
	reader.ReadU16()
	reader.Read(4)
	reader.ReadU8()
	reader.ReadU8()
	p.At = util.FromFileTime(reader.ReadU64())
	reader.ReadU8()
	reader.ReadU8()
	reader.ReadU8()
	p.Today = reader.ReadU16()
	p.Days = reader.ReadU8()
}
