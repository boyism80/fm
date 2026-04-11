package response

import (
	"github.com/boyism80/fm/stream"
)

type NpcAction struct {
	Bytes []byte
}

func (p *NpcAction) Opcode() uint16 {
	return 0xBE
}

func (p *NpcAction) Serialize(writer *stream.StreamWriter) error {
	writer.Write(p.Bytes)
	return nil
}

func (s *NpcAction) Deserialize(reader *stream.StreamReader) {
}
