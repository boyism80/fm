package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

// GiveDebuff notifies the client that a debuff was applied to self (opcode 0x15).
type GiveDebuff struct {
	Disease    constant.DebuffFlag
	X          int16
	SkillID    uint16
	SkillLevel uint16
	DurationMs int32
}

func (p *GiveDebuff) Opcode() uint16 { return 0x15 }

func (p *GiveDebuff) Serialize(writer *stream.StreamWriter) error {
	if err := WriteSingleMask(writer, SlotFromDebuffFlag(p.Disease)); err != nil {
		return err
	}
	writer.Write16(p.X)
	writer.Write16(int16(p.SkillID))
	writer.Write16(int16(p.SkillLevel))
	writer.Write32(p.DurationMs)
	writer.Write16(0)
	writer.Write16(0)
	writer.WriteU8(1)
	return nil
}

func (p *GiveDebuff) Deserialize(_ *stream.StreamReader) error { return nil }
