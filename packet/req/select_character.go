package req

import (
	"github.com/boyism80/fm/stream"
)

type SelectCharacter struct {
	CharacterId uint32
}

func (a *SelectCharacter) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *SelectCharacter) Deserialize(reader *stream.StreamReader) error {
	id, err := reader.ReadU32()
	if err != nil {
		return err
	}

	a.CharacterId = id
	return nil
}
