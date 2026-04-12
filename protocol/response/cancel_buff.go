package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

type CancelBuff struct {
	Buffs []constant.BuffFlag
}

func (p *CancelBuff) Opcode() uint16 { return 0x16 }

func (p *CancelBuff) Serialize(writer *stream.StreamWriter) error {
	if p.Buffs == nil {
		p.Buffs = []constant.BuffFlag{}
	}
	WriteBuffs(writer, p.Buffs)
	writer.WriteU8(1)
	return nil
}

func (p *CancelBuff) Deserialize(_ *stream.StreamReader) {}
