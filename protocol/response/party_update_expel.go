package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyUpdateExpel struct {
	ForChannel          int32
	PartyID             uint32
	TargetCharacterID   uint32
	TargetCharacterName string
	LeaderCharacterID   uint32
	Members             []PartyMemberStatus
}

func (p *PartyUpdateExpel) Opcode() uint16 {
	return 0x2D
}

func (p *PartyUpdateExpel) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.PartyS2CPartyUpdate))
	w.WriteU32(p.PartyID)
	w.WriteU32(p.TargetCharacterID)
	w.WriteU8(1)
	w.WriteU8(1)
	w.WriteStr16(p.TargetCharacterName)
	writePartyStatusBlock(w, p.ForChannel, p.LeaderCharacterID, p.Members, false)
	return nil
}

func (p *PartyUpdateExpel) Deserialize(*stream.StreamReader) {}
