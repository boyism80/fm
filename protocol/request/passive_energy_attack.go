package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type PassiveEnergyAttack struct {
	dto.CloseAttackInfo
}

func (*PassiveEnergyAttack) Opcode() byte { return 0x1E }

func (p *PassiveEnergyAttack) Serialize(writer *stream.StreamWriter) error {
	return p.CloseAttackInfo.Serialize(writer)
}

func (p *PassiveEnergyAttack) Deserialize(reader *stream.StreamReader) {
	p.CloseAttackInfo.Deserialize(reader)
}
