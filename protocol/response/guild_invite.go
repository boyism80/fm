package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildInvite struct {
	GuildID     uint32
	InviterName string
}

func (p *GuildInvite) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildInvite) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CInvite))
	w.WriteU32(p.GuildID)
	w.WriteStr16(p.InviterName)
	return nil
}

func (p *GuildInvite) Deserialize(*stream.StreamReader) {}
