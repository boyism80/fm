package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildMemberLevelClassUpdate struct {
	GuildID     uint32
	CharacterID uint32
	Level       uint32
	ClassID     uint32
}

func (p *GuildMemberLevelClassUpdate) Opcode() uint16 {
	return 0x30
}

func (p *GuildMemberLevelClassUpdate) Serialize(w *stream.StreamWriter) error {
	w.WriteU8(uint8(constant.GuildS2CMemberLevelClassUpdate))
	w.WriteU32(p.GuildID)
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.Level)
	w.WriteU32(p.ClassID)
	return nil
}

func (p *GuildMemberLevelClassUpdate) Deserialize(*stream.StreamReader) {}
