package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AllianceChangeGuildMembers struct {
	Added      bool
	AllianceID uint32
	GuildID    uint32
	Members    []dto.AllianceGuildMemberRank
}

func (p *AllianceChangeGuildMembers) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceChangeGuildMembers) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CChangeGuildMembers))
	if p.Added {
		w.WriteU32(p.AllianceID)
	} else {
		w.WriteU32(0)
	}
	w.WriteU32(p.GuildID)
	writeAllianceGuildMemberRanks(w, p.Added, p.Members)
	return nil
}

func (p *AllianceChangeGuildMembers) Deserialize(*stream.StreamReader) {}
