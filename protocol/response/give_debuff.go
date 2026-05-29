package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type GiveDebuff struct {
	CharacterID int32
	Debuff      constant.DebuffFlag
	X           int16
	SkillID     uint16
	SkillLevel  uint16
	DurationMs  int32
}

func (p *GiveDebuff) Opcode() uint16 { return 0x90 }

func (p *GiveDebuff) Serialize(writer *stream.StreamWriter) error {
	writer.Write32(p.CharacterID)
	WriteDebuff(writer, p.Debuff)
	if p.Debuff == constant.DebuffFlagPoison {
		writer.Write16(p.X)
	}
	writer.Write16(int16(p.SkillID))
	writer.Write16(int16(p.SkillLevel))
	writer.Write16(0)
	writer.Write16(0)
	writer.WriteU8(1)
	return nil
}

func (p *GiveDebuff) Deserialize(_ *stream.StreamReader) {}
