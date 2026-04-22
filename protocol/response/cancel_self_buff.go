package response

import (
	"github.com/boyism80/fm/services/game/constant"
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

func (p *CancelSelfBuff) Deserialize(_ *stream.StreamReader) {}
