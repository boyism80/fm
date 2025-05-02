package movement

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type ChairMovement struct {
	Command  uint8
	Position types.Vec2[int16]
	Unknown  int16
	NewState uint8
	Duration int16
}

func (m *ChairMovement) GetNewState() uint8 {
	return m.NewState
}

func (m *ChairMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Unknown)
	writer.WriteU8(m.NewState)
	writer.Write16(m.Duration)
}
