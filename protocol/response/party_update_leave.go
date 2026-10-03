package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type PartyUpdateLeave struct {
	ForChannel          int32
	PartyID             uint32
	TargetCharacterID   uint32
	TargetCharacterName string
	LeaderCharacterID   uint32
	Members             []PartyMemberStatus
}

func (p *PartyUpdateLeave) Opcode() uint16 {
	return 0x2D
}

func (p *PartyUpdateLeave) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.PartyS2CPartyUpdate))
	w.WriteU32(p.PartyID)
	w.WriteU32(p.TargetCharacterID)
	w.WriteU8(1)
	w.WriteU8(0)
	w.WriteStr16(p.TargetCharacterName)
	writePartyStatusBlock(w, p.ForChannel, p.LeaderCharacterID, p.Members, true)
	return nil
}

func (p *PartyUpdateLeave) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.PartyID = reader.ReadU32()
	p.TargetCharacterID = reader.ReadU32()
	reader.ReadU8()
	reader.ReadU8()
	p.TargetCharacterName = reader.ReadStr16()
	p.LeaderCharacterID, p.Members = readPartyStatusBlock(reader)
}
