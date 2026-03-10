package request

import (
	"github.com/boyism80/fm/stream"
)

// Unknown0C represents an unknown packet (0x0C)
type Unknown0C struct {
}

func (p *Unknown0C) Opcode() uint16 {
	return 0x0C
}

func (p *Unknown0C) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Unknown0C) Deserialize(reader *stream.StreamReader) error {
	return nil
}
