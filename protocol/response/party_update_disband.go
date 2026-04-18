package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyUpdateDisband struct {
	PartyID           uint32
	LeaderCharacterID uint32
}

func (p *PartyUpdateDisband) Opcode() uint16 {
	return 0x2D
}

func (p *PartyUpdateDisband) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.PartyS2CPartyUpdate))
	w.WriteU32(p.PartyID)
	w.WriteU32(p.LeaderCharacterID)
	w.WriteU8(0)
	w.WriteU32(p.LeaderCharacterID)
	return nil
}

func (p *PartyUpdateDisband) Deserialize(*stream.StreamReader) {}
