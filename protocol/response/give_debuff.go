package response

import (
	"github.com/boyism80/fm/maple/constant"
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

func (p *GiveDebuff) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.Read32()
	for i := 0; i < constant.MaxBuffFlag; i++ {
		if v := reader.Read32(); v != 0 {
			p.Debuff = constant.DebuffFlag{Mask: uint32(v), Position: i + 1}
		}
	}
	if p.Debuff.Mask == constant.DebuffFlagPoison.Mask && p.Debuff.Position == constant.DebuffFlagPoison.Position {
		p.X = reader.Read16()
	}
	p.SkillID = uint16(reader.Read16())
	p.SkillLevel = uint16(reader.Read16())
}
