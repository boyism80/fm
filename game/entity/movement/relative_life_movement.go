package movement

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type RelativeLifeMovement struct {
	Command  byte
	Position types.Vec2[int16]
	NewState uint8
	Duration int16
}

func (m *RelativeLifeMovement) GetNewState() uint8 {
	return m.NewState
}

func (m *RelativeLifeMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.WriteU8(m.NewState)
	writer.Write16(m.Duration)
}
