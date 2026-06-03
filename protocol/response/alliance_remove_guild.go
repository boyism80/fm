package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AllianceRemoveGuild struct {
	Info           *dto.AllianceInfo
	RemovedGuildID uint32
	Guild          *dto.GuildInfo
	Expelled       bool
}

func (p *AllianceRemoveGuild) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceRemoveGuild) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CRemoveGuild))
	if p.Info != nil {
		writeAllianceInfo(w, p.Info)
	}
	w.WriteU32(p.RemovedGuildID)
	if p.Guild != nil {
		writeGuildInfo(w, p.Guild)
	}
	w.WriteBoolean(p.Expelled)
	return nil
}

func (p *AllianceRemoveGuild) Deserialize(*stream.StreamReader) {}
