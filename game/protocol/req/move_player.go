package req

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/protocol"
)

type MovePlayer struct {
	Fragments []protocol.MoveFragment
}

func (m *MovePlayer) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (m *MovePlayer) Deserialize(reader *stream.StreamReader) error {
	if _, err := reader.Read(5); err != nil {
		return err
	}

	fragments, err := protocol.Parse(reader)
	if err != nil {
		return err
	}

	m.Fragments = fragments
	return nil
}
