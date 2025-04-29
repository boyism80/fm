package req

import (
	"github.com/boyism80/fm/stream"
)

type LoginGame struct {
	PlayerId uint32
}

func (a *LoginGame) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *LoginGame) Deserialize(reader *stream.StreamReader) error {
	id, err := reader.ReadU32()
	if err != nil {
		return err
	}

	a.PlayerId = id
	return nil
}
