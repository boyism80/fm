package response

import (
	"github.com/boyism80/fm/stream"
)

type EndOfServerList struct{}

func (e *EndOfServerList) Opcode() uint16 {
	return 0x02
}

func (e *EndOfServerList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0xFF)

	return nil
}

func (e *EndOfServerList) Deserialize(reader *stream.StreamReader) {
}
