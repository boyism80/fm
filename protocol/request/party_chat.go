package request

import (
	"github.com/boyism80/fm/stream"
)

type PartyChat struct {
	Type          byte
	NumRecipients byte
	Recipients    []int32
	Message       string
}

func (m *PartyChat) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *PartyChat) Deserialize(reader *stream.StreamReader) {
	m.Type = reader.ReadU8()
	m.NumRecipients = reader.ReadU8()
	if m.NumRecipients <= 0 {
		return
	}
	n := int(m.NumRecipients)
	m.Recipients = make([]int32, n)
	for i := 0; i < n; i++ {
		m.Recipients[i] = reader.Read32()
	}
	m.Message = reader.ReadStr16()
}
