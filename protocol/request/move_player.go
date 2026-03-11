package request

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MovePlayer struct {
	Fragments []dto.MoveFragment
}

func (m *MovePlayer) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MovePlayer) Deserialize(reader *stream.StreamReader) error {
	if _, err := reader.Read(5); err != nil {
		return err
	}

	fragments, err := dto.ReadMovements(reader)
	if err != nil {
		return err
	}

	m.Fragments = fragments
	return nil
}
