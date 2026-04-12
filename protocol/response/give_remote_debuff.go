package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

type GiveRemoteDebuff struct {
	CharacterID int32
	Disease     constant.DebuffFlag
	X           int16
	SkillID     uint16
	SkillLevel  uint16
}

func (p *GiveRemoteDebuff) Opcode() uint16 { return 0x90 }

func (p *GiveRemoteDebuff) Serialize(writer *stream.StreamWriter) error {
	writer.Write32(p.CharacterID)
	WriteDebuff(writer, p.Disease)
	if p.Disease == constant.DebuffFlagPoison {
		writer.Write16(p.X)
	}
	writer.Write16(int16(p.SkillID))
	writer.Write16(int16(p.SkillLevel))
	writer.Write16(0)
	writer.Write16(0)
	writer.WriteU8(1)
	return nil
}

func (p *GiveRemoteDebuff) Deserialize(_ *stream.StreamReader) {}
