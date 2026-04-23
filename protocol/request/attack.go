package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type Attack struct {
	dto.CloseAttackInfo
}

func (*Attack) Opcode() byte { return 0x1B }

func (a *Attack) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *Attack) Deserialize(reader *stream.StreamReader) {
	a.CloseAttackInfo.Deserialize(reader)
}
