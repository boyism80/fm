package req

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity/movement"
)

type MovePlayer struct {
	Fragments []movement.MoveFragment
}

func (m *MovePlayer) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MovePlayer) Deserialize(reader *stream.StreamReader) error {
	if _, err := reader.Read(5); err != nil { // skip
		return err
	}

	fragments, err := movement.Parse(reader)
	if err != nil {
		return err
	}

	m.Fragments = fragments
	return nil
}
