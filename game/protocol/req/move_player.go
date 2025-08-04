package req

import (
	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/game/action"
)

type MovePlayer struct {
	Fragments []action.MoveFragment
}

func (m *MovePlayer) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MovePlayer) Deserialize(reader *stream.StreamReader) error {
	if _, err := reader.Read(5); err != nil {
		return err
	}

	fragments, err := action.ReadMovements(reader)
	if err != nil {
		return err
	}

	m.Fragments = fragments
	return nil
}
