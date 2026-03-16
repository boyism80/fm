package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type SummonAttack struct {
	dto.AttackInfo
}

func (a *SummonAttack) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *SummonAttack) Deserialize(reader *stream.StreamReader) error {
	return a.AttackInfo.Deserialize(reader, 0x8D)
}
