package req

import (
	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/game/action"
)

type Attack struct {
	action.AttackInfo
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
