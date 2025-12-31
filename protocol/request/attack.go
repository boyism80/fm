package request

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/protocol/dto"
)

type Attack struct {
	dto.AttackInfo
}

func (a *Attack) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *Attack) Deserialize(reader *stream.StreamReader) error {
	err := a.AttackInfo.Deserialize(reader, 0x1B)
	if err != nil {
		return err
	}
	return nil
}
