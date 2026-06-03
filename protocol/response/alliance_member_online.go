package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type AllianceMemberOnline struct {
	AllianceID  uint32
	GuildID     uint32
	CharacterID uint32
	Online      bool
}

func (p *AllianceMemberOnline) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceMemberOnline) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CMemberOnline))
	w.WriteU32(p.AllianceID)
	w.WriteU32(p.GuildID)
	w.WriteU32(p.CharacterID)
	w.WriteBoolean(p.Online)
	return nil
}

func (p *AllianceMemberOnline) Deserialize(*stream.StreamReader) {}
