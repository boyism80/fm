package response

import (
	"github.com/boyism80/fm/stream"
)

type StopClock struct {
}

func (p *StopClock) Opcode() uint16 {
	return 0x6B
}

func (p *StopClock) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *StopClock) Deserialize(reader *stream.StreamReader) {
}
