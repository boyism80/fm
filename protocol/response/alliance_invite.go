package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type AllianceInvite struct {
	InviterGuildID uint32
	InviterName    string
	AllianceName   string
}

func (p *AllianceInvite) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceInvite) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CInvite))
	w.WriteU32(p.InviterGuildID)
	w.WriteStr16(p.InviterName)
	w.WriteStr16(p.AllianceName)
	return nil
}

func (p *AllianceInvite) Deserialize(*stream.StreamReader) {}
