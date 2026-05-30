package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildUpdateGP struct {
	GuildID    uint32
	GP         uint32
	GuildLevel uint32
}

func (p *GuildUpdateGP) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildUpdateGP) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CUpdateGP))
	w.WriteU32(p.GuildID)
	w.WriteU32(p.GP)
	w.WriteU32(p.GuildLevel)
	return nil
}

func (p *GuildUpdateGP) Deserialize(*stream.StreamReader) {}
