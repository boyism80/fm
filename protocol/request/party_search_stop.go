package request

import (
	"github.com/boyism80/fm/stream"
)

type PartySearchStop struct {
}

func (p *PartySearchStop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *PartySearchStop) Deserialize(reader *stream.StreamReader) {
}
