package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyStatusMessage struct {
	Code constant.PartyStatusCode
	Name string
}

func (p *PartyStatusMessage) Opcode() uint16 {
	return 0x2D
}

func (p *PartyStatusMessage) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(p.Code))
	if p.Name != "" {
		w.WriteStr16(p.Name)
	}
	return nil
}

func (p *PartyStatusMessage) Deserialize(reader *stream.StreamReader) {
	p.Code = constant.PartyStatusCode(reader.ReadU8())
	if reader.Remaining() > 0 {
		p.Name = reader.ReadStr16()
	}
}
