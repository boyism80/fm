package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyInvite struct {
	PartyID     uint32
	InviterName string
	PartySearch bool
}

func (p *PartyInvite) Opcode() uint16 {
	return 0x2D
}

func (p *PartyInvite) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.PartyS2CInvite))
	w.WriteU32(p.PartyID)
	w.WriteStr16(p.InviterName)
	w.WriteBoolean(p.PartySearch)
	return nil
}

func (p *PartyInvite) Deserialize(*stream.StreamReader) {}
