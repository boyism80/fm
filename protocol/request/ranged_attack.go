package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type RangedAttack struct {
	dto.AttackInfo
}

func (r *RangedAttack) Serialize(writer *stream.StreamWriter) error {
	// TODO: implement after packet structure is analyzed
	return nil
}

func (r *RangedAttack) Deserialize(reader *stream.StreamReader) error {
	// TODO: implement after packet structure is analyzed
	return nil
}
