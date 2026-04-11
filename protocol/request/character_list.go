package request

import (
	"github.com/boyism80/fm/stream"
)

type CharacterList struct {
	Server  uint8
	Channel uint8
}

func (a *CharacterList) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *CharacterList) Deserialize(reader *stream.StreamReader) {
	reader.Read(1)
	a.Server = reader.ReadU8()
	a.Channel = reader.ReadU8()
}
