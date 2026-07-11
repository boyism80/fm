package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/boyism80/fm/stream"
)

type QuestStatus struct {
	QuestID        uint16
	Status         uint8
	MobKills       []uint16
	Deadline       time.Time
	StatusRecord   string
	CompletionTime time.Time
}

func (q *QuestStatus) WriteStartedPayload(writer *stream.StreamWriter) {
	if q == nil {
		writer.Write([]byte{0x00, 0x00})
		return
	}
	if len(q.MobKills) > 0 {
		var sb strings.Builder
		for _, kills := range q.MobKills {
			sb.WriteString(fmt.Sprintf("%03d", kills))
		}
		writer.WriteStr16(sb.String())
	} else if !q.Deadline.IsZero() {
		writer.WriteU16(9)
		writer.WriteU8(1)
		writer.WriteDateTime(q.Deadline)
	} else if q.StatusRecord != "" {
		writer.WriteStr16(q.StatusRecord)
	} else {
		writer.Write([]byte{0x00, 0x00})
	}
}
