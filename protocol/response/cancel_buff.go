package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

// CancelBuff notifies the client to remove buff(s) from self (buff off).
type CancelBuff struct {
	Buffs []constant.BuffFlag
}

func (p *CancelBuff) Opcode() uint16 { return 0x16 }

func (p *CancelBuff) Serialize(writer *stream.StreamWriter) error {
	if p.Buffs == nil {
		p.Buffs = []constant.BuffFlag{}
	}
	if err := WriteMask(writer, SlotsFromBuffFlags(p.Buffs)); err != nil {
		return err
	}
	writer.WriteU8(1)
	return nil
}

func (p *CancelBuff) Deserialize(_ *stream.StreamReader) error { return nil }
