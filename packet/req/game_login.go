package req

import (
	"github.com/boyism80/fm/stream"
)

type GameLogin struct {
	CharacterId uint32
}

func (a *GameLogin) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *GameLogin) Deserialize(reader *stream.StreamReader) error {
	id, err := reader.ReadU32()
	if err != nil {
		return err
	}

	a.CharacterId = id
	return nil
}
