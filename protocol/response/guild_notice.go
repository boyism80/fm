package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildNotice struct {
	GuildID uint32
	Notice  string
}

func (p *GuildNotice) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildNotice) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CNotice))
	w.WriteU32(p.GuildID)
	w.WriteStr16(p.Notice)
	return nil
}

func (p *GuildNotice) Deserialize(*stream.StreamReader) {}
