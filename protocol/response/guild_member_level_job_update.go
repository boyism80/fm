package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildMemberLevelJobUpdate struct {
	GuildID     uint32
	CharacterID uint32
	Level       uint32
	JobID       uint32
}

func (p *GuildMemberLevelJobUpdate) Opcode() uint16 {
	return 0x30
}

func (p *GuildMemberLevelJobUpdate) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CMemberLevelJobUpdate))
	w.WriteU32(p.GuildID)
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.Level)
	w.WriteU32(p.JobID)
	return nil
}

func (p *GuildMemberLevelJobUpdate) Deserialize(*stream.StreamReader) {}
