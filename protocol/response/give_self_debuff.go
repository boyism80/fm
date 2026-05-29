package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type GiveSelfDebuff struct {
	Debuff     constant.DebuffFlag
	X          int16
	SkillID    uint16
	SkillLevel uint16
	DurationMs int32
}

func (p *GiveSelfDebuff) Opcode() uint16 { return 0x15 }

func (p *GiveSelfDebuff) Serialize(writer *stream.StreamWriter) error {
	WriteDebuff(writer, p.Debuff)
	writer.Write16(p.X)
	writer.Write16(int16(p.SkillID))
	writer.Write16(int16(p.SkillLevel))
	writer.Write32(p.DurationMs)
	writer.Write16(0)
	writer.Write16(0)
	writer.WriteU8(1)
	return nil
}

func (p *GiveSelfDebuff) Deserialize(_ *stream.StreamReader) {}
