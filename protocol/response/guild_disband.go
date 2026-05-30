package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildDisband struct {
	GuildID uint32
}

func (p *GuildDisband) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildDisband) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CDisband))
	w.WriteU32(p.GuildID)
	w.WriteU8(1)
	return nil
}

func (p *GuildDisband) Deserialize(*stream.StreamReader) {}
