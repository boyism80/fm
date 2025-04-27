package resp

import (
	"github.com/boyism80/fm/stream"
)

type EndOfServerList struct{}

func (e *EndOfServerList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x02)
	writer.WriteU8(0xFF)

	return nil
}

func (e *EndOfServerList) Deserialize(reader *stream.StreamReader) error {
	return nil
}
