package response

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/stream"
)

type CancelSelfBuff struct {
	Buffs []constant.BuffFlag
}

func (p *CancelSelfBuff) Opcode() uint16 { return 0x16 }

func (p *CancelSelfBuff) Serialize(writer *stream.StreamWriter) error {
	if p.Buffs == nil {
		p.Buffs = []constant.BuffFlag{}
	}
	WriteBuffs(writer, p.Buffs)
	writer.WriteU8(1)
	return nil
}

func (p *CancelSelfBuff) Deserialize(reader *stream.StreamReader) {
	for i := 0; i < constant.MaxBuffFlag; i++ {
		mask := uint32(reader.Read32())
		for bit := uint32(1); bit != 0; bit <<= 1 {
			if mask&bit != 0 {
				p.Buffs = append(p.Buffs, constant.BuffFlag{Mask: bit, Position: i + 1})
			}
		}
	}
}
