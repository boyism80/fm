package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type CancelBuff struct {
	CharacterID int32
	Buffs       []constant.BuffFlag
}

func (p *CancelBuff) Opcode() uint16 { return 0x91 }

func (p *CancelBuff) Serialize(writer *stream.StreamWriter) error {
	if p.Buffs == nil {
		p.Buffs = []constant.BuffFlag{}
	}
	writer.Write32(p.CharacterID)
	WriteBuffs(writer, p.Buffs)
	return nil
}

func (p *CancelBuff) Deserialize(_ *stream.StreamReader) {}
