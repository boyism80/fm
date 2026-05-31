package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildMemberOnline struct {
	GuildID     uint32
	CharacterID uint32
	Online      bool
}

func (p *GuildMemberOnline) Opcode() uint16 {
	return 0x30
}

func (p *GuildMemberOnline) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CMemberOnline))
	w.WriteU32(p.GuildID)
	w.WriteU32(p.CharacterID)
	online := uint8(0)
	if p.Online {
		online = 1
	}
	w.WriteU8(online)
	return nil
}

func (p *GuildMemberOnline) Deserialize(*stream.StreamReader) {}
