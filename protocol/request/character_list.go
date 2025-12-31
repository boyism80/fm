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

func (a *CharacterList) Deserialize(reader *stream.StreamReader) error {
	reader.Read(1)
	server, err := reader.ReadU8()
	if err != nil {
		return err
	}
	a.Server = server
	channel, err := reader.ReadU8()
	if err != nil {
		return err
	}
	a.Channel = channel
	return nil
}
