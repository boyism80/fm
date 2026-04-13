package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type GiveHomingBeacon struct {
	SkillID uint32
	MobOID  uint32
	X       int32
}

func (p *GiveHomingBeacon) Opcode() uint16 { return 0x15 }

func (p *GiveHomingBeacon) Serialize(writer *stream.StreamWriter) error {
	WriteBuffs(writer, []constant.BuffFlag{constant.BuffFlagHomingBeacon})
	writer.WriteU16(0)
	writer.Write32(p.X)
	writer.Write64(int64(p.SkillID))
	writer.WriteU8(0)
	writer.Write64(int64(p.MobOID))
	writer.WriteU16(0)
	writer.WriteU16(0)
	writer.WriteU8(0)
	return nil
}

func (p *GiveHomingBeacon) Deserialize(_ *stream.StreamReader) {}

type CancelHomingBeacon struct{}

func (p *CancelHomingBeacon) Opcode() uint16 { return 0x16 }

func (p *CancelHomingBeacon) Serialize(writer *stream.StreamWriter) error {
	WriteBuffs(writer, []constant.BuffFlag{constant.BuffFlagHomingBeacon})
	return nil
}

func (p *CancelHomingBeacon) Deserialize(_ *stream.StreamReader) {}
