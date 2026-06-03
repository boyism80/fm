package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AllianceCreate struct {
	Info   *dto.AllianceInfo
	Guilds []*dto.GuildInfo
}

func (p *AllianceCreate) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceCreate) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CCreate))
	if p.Info != nil {
		writeAllianceInfo(w, p.Info)
	}
	for _, g := range p.Guilds {
		if g == nil {
			return nil
		}
		writeGuildInfo(w, g)
	}
	return nil
}

func (p *AllianceCreate) Deserialize(*stream.StreamReader) {}
