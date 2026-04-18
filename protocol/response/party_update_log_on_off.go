package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyUpdateLogOnOff struct {
	ForChannel        int32
	PartyID           uint32
	LeaderCharacterID uint32
	Members           []PartyMemberStatus
}

func (p *PartyUpdateLogOnOff) Opcode() uint16 {
	return 0x2D
}

func (p *PartyUpdateLogOnOff) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.PartyS2CSilentUpdate))
	w.WriteU32(p.PartyID)
	writePartyStatusBlock(w, p.ForChannel, p.LeaderCharacterID, p.Members, true)
	return nil
}

func (p *PartyUpdateLogOnOff) Deserialize(*stream.StreamReader) {}
