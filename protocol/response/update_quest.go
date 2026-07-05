package response

import (
	"strconv"
	"strings"
	"time"

	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

const (
	QuestWireStatusNotStarted uint8 = 0
	QuestWireStatusStarted    uint8 = 1
	QuestWireStatusCompleted  uint8 = 2
)

type UpdateQuest struct {
	QuestID        uint16
	Status         uint8
	StartedPayload string
	CompletionTime time.Time
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
		if p.StartedPayload != "" {
			if strings.HasPrefix(p.StartedPayload, "time_") {
				writer.WriteU16(9)
				writer.WriteU8(1)
				timeVal, err := strconv.ParseInt(p.StartedPayload[5:], 10, 64)
				if err != nil {
					writer.WriteU64(util.ToFileTime(time.Time{}))
				} else {
					writer.WriteU64(util.ToFileTime(util.GetTime(timeVal)))
				}
			} else {
				writer.WriteStr8(p.StartedPayload)
			}
		} else {
			writer.Write([]byte{0x00, 0x00})
		}

	case QuestWireStatusCompleted:
		completionTime := p.CompletionTime
		if completionTime.IsZero() {
			completionTime = time.Now()
		}
		writer.WriteU64(util.ToFileTime(completionTime))
		writer.WriteU16(0)
	}

	return nil
}

func (p *UpdateQuest) Deserialize(reader *stream.StreamReader) {
}
