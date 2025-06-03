package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type NpcAction struct {
	Bytes []byte
}

func (p *NpcAction) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0xBE)
	writer.Write(p.Bytes)
	return nil
}

func (s *NpcAction) Deserialize(reader *stream.StreamReader) error {
	return nil
}
