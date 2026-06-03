package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AllianceChangeMembership struct {
	InAlliance bool
	AllianceID uint32
	Guilds     []dto.AllianceMembershipChangeGuild
}

func (p *AllianceChangeMembership) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceChangeMembership) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CChangeMembership))
	w.WriteBoolean(p.InAlliance)
	if p.InAlliance {
		w.WriteU32(p.AllianceID)
	} else {
		w.WriteU32(0)
	}
	writeAllianceMembershipChangeGuilds(w, p.InAlliance, p.Guilds)
	return nil
}

func (p *AllianceChangeMembership) Deserialize(*stream.StreamReader) {}
