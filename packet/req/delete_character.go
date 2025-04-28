package req

import "github.com/boyism80/fm/stream"

type DeleteCharacter struct {
	Id uint32
}

func (a *DeleteCharacter) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *DeleteCharacter) Deserialize(reader *stream.StreamReader) error {
	_, err := reader.Read(5)
	if err != nil {
		return err
	}
	id, err := reader.ReadU32()
	if err != nil {
		return err
	}

	a.Id = id
	return nil
}
