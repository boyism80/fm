package request

import "github.com/boyism80/fm/stream"

type DeleteCharacter struct {
	ID uint32
}

func (a *DeleteCharacter) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *DeleteCharacter) Deserialize(reader *stream.StreamReader) {
	reader.Read(5)
	a.ID = reader.ReadU32()
}
