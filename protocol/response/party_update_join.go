package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyUpdateJoin struct {
	ForChannel           int32
	PartyID              uint32
	JoiningCharacterName string
	LeaderCharacterID    uint32
	Members              []PartyMemberStatus
}

func (p *PartyUpdateJoin) Opcode() uint16 {
	return 0x2D
}

func (p *PartyUpdateJoin) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.PartyS2CPartyJoin))
	w.WriteU32(p.PartyID)
	w.WriteStr16(p.JoiningCharacterName)
	writePartyStatusBlock(w, p.ForChannel, p.LeaderCharacterID, p.Members, false)
	return nil
}

func (p *PartyUpdateJoin) Deserialize(*stream.StreamReader) {}
