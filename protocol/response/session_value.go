package response

import (
	"github.com/boyism80/fm/stream"
)

type SessionValue struct {
	Key   string
	Value string
}

func (p *SessionValue) Opcode() uint16 {
	return 0x47
}

func (p *SessionValue) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(p.Key)
	writer.WriteStr16(p.Value)
	return nil
}

func (p *SessionValue) Deserialize(reader *stream.StreamReader) {
	p.Key = reader.ReadStr16()
	p.Value = reader.ReadStr16()
}
