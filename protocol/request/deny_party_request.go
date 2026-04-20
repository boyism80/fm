package request

import "github.com/boyism80/fm/stream"

type DenyPartyRequest struct {
	Action      uint8
	InviterName string
}

func (p *DenyPartyRequest) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DenyPartyRequest) Deserialize(reader *stream.StreamReader) {
	p.Action = reader.ReadU8()
	p.InviterName = reader.ReadStr16()
}
