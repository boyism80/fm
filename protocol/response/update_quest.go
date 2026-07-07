package response

import (
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

const (
	QuestWireStatusNotStarted uint8 = 0
	QuestWireStatusStarted    uint8 = 1
	QuestWireStatusCompleted  uint8 = 2
)

type UpdateQuest struct {
	dto.QuestStatus
}

func (p *UpdateQuest) Opcode() uint16 {
	return 0x1C
}

func (p *UpdateQuest) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)
	writer.WriteU16(p.QuestID)
	writer.WriteU8(p.Status)

	switch p.Status {
	case QuestWireStatusNotStarted:
		writer.Write(make([]byte, 10))

	case QuestWireStatusStarted:
		p.WriteStartedPayload(writer)

	case QuestWireStatusCompleted:
		completionTime := p.CompletionTime
		if completionTime.IsZero() {
			completionTime = clock.Now()
		}
		writer.WriteU64(util.ToFileTime(completionTime))
		writer.WriteU16(0)
	}

	return nil
}

func (p *UpdateQuest) Deserialize(reader *stream.StreamReader) {
}
