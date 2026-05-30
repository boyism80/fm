package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildRankTitleChange struct {
	GuildID    uint32
	RankTitles [5]string
}

func (p *GuildRankTitleChange) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildRankTitleChange) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CRankTitleChange))
	w.WriteU32(p.GuildID)
	for i := 0; i < 5; i++ {
		w.WriteStr16(p.RankTitles[i])
	}
	return nil
}

func (p *GuildRankTitleChange) Deserialize(*stream.StreamReader) {}
