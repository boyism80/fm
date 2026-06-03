package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AllianceAddGuild struct {
	Info       *dto.AllianceInfo
	NewGuildID uint32
	Guild      *dto.GuildInfo
}

func (p *AllianceAddGuild) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceAddGuild) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CAddGuild))
	if p.Info != nil {
		writeAllianceInfo(w, p.Info)
	}
	w.WriteU32(p.NewGuildID)
	if p.Guild != nil {
		writeGuildInfo(w, p.Guild)
	}
	w.WriteU8(0)
	return nil
}

func (p *AllianceAddGuild) Deserialize(*stream.StreamReader) {}
