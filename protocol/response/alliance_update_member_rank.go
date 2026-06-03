package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type AllianceUpdateMemberRank struct {
	AllianceID   uint32
	CharacterID  uint32
	AllianceRank uint32
}

func (p *AllianceUpdateMemberRank) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceUpdateMemberRank) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CUpdateMemberRank))
	w.WriteU32(p.AllianceID)
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.AllianceRank)
	return nil
}

func (p *AllianceUpdateMemberRank) Deserialize(*stream.StreamReader) {}
