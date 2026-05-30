package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildEmblemChange struct {
	GuildID     uint32
	LogoBG      uint16
	LogoBGColor uint8
	Logo        uint16
	LogoColor   uint8
}

func (p *GuildEmblemChange) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildEmblemChange) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CEmblemChange))
	w.WriteU32(p.GuildID)
	w.WriteU16(p.LogoBG)
	w.WriteU8(p.LogoBGColor)
	w.WriteU16(p.Logo)
	w.WriteU8(p.LogoColor)
	return nil
}

func (p *GuildEmblemChange) Deserialize(*stream.StreamReader) {}
