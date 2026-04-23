package request

import (
	"github.com/boyism80/fm/stream"
)

type PartySearchStop struct {
}

func (*PartySearchStop) Opcode() byte { return 0xB6 }

func (p *PartySearchStop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *PartySearchStop) Deserialize(reader *stream.StreamReader) {
}
