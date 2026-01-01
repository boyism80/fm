package request

import (
	"github.com/boyism80/fm/stream"
)

// PartySearchStop represents the PARTY_SEARCH_STOP packet (0xB6)
type PartySearchStop struct {
}

func (p *PartySearchStop) Opcode() uint16 {
	return 0xB6
}

func (p *PartySearchStop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *PartySearchStop) Deserialize(reader *stream.StreamReader) error {
	return nil
}

