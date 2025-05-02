package movement

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type AbsoluteLifeMovement struct {
	Command        uint8
	Position       types.Vec2[int16]
	Duration       int16
	NewState       uint8
	PixelPerSecond types.Vec2[int16]
	Offset         types.Vec2[int16]
	Unknown        int16
}

func (m *AbsoluteLifeMovement) GetNewState() uint8 {
	return m.NewState
}

func (m *AbsoluteLifeMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.PixelPerSecond.X)
	writer.Write16(m.PixelPerSecond.Y)
	writer.Write16(m.Unknown)
	writer.WriteU8(m.NewState)
	writer.Write16(m.Duration)
}
