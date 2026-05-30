package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type GuildShowRanks struct {
	NPCID   uint32
	Entries []dto.GuildRankingEntry
}

func (p *GuildShowRanks) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildShowRanks) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CShowRanks))
	w.WriteU32(p.NPCID)
	w.WriteU32(uint32(len(p.Entries)))
	for _, entry := range p.Entries {
		w.WriteStr16(entry.Name)
		w.WriteU32(entry.GP)
		w.WriteU32(entry.Logo)
		w.WriteU32(entry.LogoColor)
		w.WriteU32(entry.LogoBG)
		w.WriteU32(entry.LogoBGColor)
	}
	return nil
}

func (p *GuildShowRanks) Deserialize(*stream.StreamReader) {}
