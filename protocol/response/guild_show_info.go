package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type GuildShowInfo struct {
	InGuild bool
	Info    *dto.GuildInfo
}

func (p *GuildShowInfo) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildShowInfo) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CShowInfo))
	if !p.InGuild || p.Info == nil {
		w.WriteU8(0)
		return nil
	}
	w.WriteU8(1)
	writeGuildInfo(w, p.Info)
	return nil
}

func (p *GuildShowInfo) Deserialize(*stream.StreamReader) {}
