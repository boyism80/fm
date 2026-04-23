package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MagicAttack struct {
	dto.MagicAttackInfo
}

func (*MagicAttack) Opcode() byte { return 0x1D }

func (m *MagicAttack) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MagicAttack) Deserialize(reader *stream.StreamReader) {
	m.MagicAttackInfo.Deserialize(reader)
}
