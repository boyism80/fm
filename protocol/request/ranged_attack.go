package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type RangedAttack struct {
	dto.RangedAttackInfo
}

func (*RangedAttack) Opcode() byte { return 0x1C }

func (r *RangedAttack) Serialize(writer *stream.StreamWriter) error {

	return nil
}

func (r *RangedAttack) Deserialize(reader *stream.StreamReader) {
	r.RangedAttackInfo.Deserialize(reader)
}
