package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AllianceShowGuilds struct {
	Guilds []*dto.GuildInfo
}

func (p *AllianceShowGuilds) Opcode() uint16 {
	return allianceSendOpcode
}

func (p *AllianceShowGuilds) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.AllianceS2CShowGuilds))
	if len(p.Guilds) == 0 {
		w.WriteU32(0)
		return nil
	}
	w.WriteU32(uint32(len(p.Guilds)))
	for _, g := range p.Guilds {
		if g == nil {
			return nil
		}
		writeGuildInfo(w, g)
	}
	return nil
}

func (p *AllianceShowGuilds) Deserialize(*stream.StreamReader) {}
