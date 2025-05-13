package req

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/protocol"
)

type Attack struct {
	protocol.AttackInfo
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
