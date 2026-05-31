package response

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type GuildMemberLeft struct {
	GuildID     uint32
	CharacterID uint32
	Name        string
	WasExpelled bool
}

func (p *GuildMemberLeft) Opcode() uint16 {
	return 0x30
}

func (p *GuildMemberLeft) Serialize(w *stream.StreamWriter) error {
	sub := constant.GuildS2CMemberLeft
	if p.WasExpelled {
		sub = constant.GuildS2CMemberExpelled
	}
	w.WriteU8(uint8(sub))
	w.WriteU32(p.GuildID)
	w.WriteU32(p.CharacterID)
	w.WriteStr16(p.Name)
	return nil
}

func (p *GuildMemberLeft) Deserialize(*stream.StreamReader) {}
