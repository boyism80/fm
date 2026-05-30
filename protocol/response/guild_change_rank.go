package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildChangeRank struct {
	GuildID     uint32
	CharacterID uint32
	GuildRank   uint8
}

func (p *GuildChangeRank) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildChangeRank) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CChangeRank))
	w.WriteU32(p.GuildID)
	w.WriteU32(p.CharacterID)
	w.WriteU8(p.GuildRank)
	return nil
}

func (p *GuildChangeRank) Deserialize(*stream.StreamReader) {}
