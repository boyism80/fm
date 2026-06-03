package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type AllianceUpdateMember struct {
	AllianceID  uint32
	GuildID     uint32
	CharacterID uint32
	Level       uint32
	ClassID     uint32
}

func (p *AllianceUpdateMember) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceUpdateMember) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CUpdateMember))
	w.WriteU32(p.AllianceID)
	w.WriteU32(p.GuildID)
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.Level)
	w.WriteU32(p.ClassID)
	return nil
}

func (p *AllianceUpdateMember) Deserialize(*stream.StreamReader) {}
