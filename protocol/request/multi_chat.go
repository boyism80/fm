package request

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type MultiChat struct {
	Type          pconst.MultiChatMode
	NumRecipients byte
	Recipients    []int32
	Message       string
}

func (*MultiChat) Opcode() byte { return 0x62 }

func (m *MultiChat) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MultiChat) Deserialize(reader *stream.StreamReader) {
	m.Type = pconst.MultiChatMode(reader.ReadU8())
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
