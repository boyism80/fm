package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type GuildNewMember struct {
	GuildID uint32
	Member  dto.GuildMemberStatus
}

func (p *GuildNewMember) Opcode() uint16 {
	return guildOperationOpcode()
}

func (p *GuildNewMember) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CNewMember))
	writeGuildMemberJoinedPayload(w, p.GuildID, p.Member)
	return nil
}

func (p *GuildNewMember) Deserialize(*stream.StreamReader) {}
