package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MagicAttack struct {
	dto.AttackInfo
}

func (m *MagicAttack) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MagicAttack) Deserialize(reader *stream.StreamReader) error {
	err := m.AttackInfo.Deserialize(reader, 0x1D)
	if err != nil {
		return err
	}
	return nil
}
