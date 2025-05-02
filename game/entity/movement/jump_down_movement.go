package movement

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type JumpDownMovement struct {
	Command        uint8
	Position       types.Vec2[int16]
	Duration       int16
	NewState       uint8
	PixelPerSecond types.Vec2[int16]
	Unknown        int16
	FH             int16
}

func (m *JumpDownMovement) GetNewState() uint8 {
	return m.NewState
}

func (m *JumpDownMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.PixelPerSecond.X)
	writer.Write16(m.PixelPerSecond.Y)
	writer.Write16(m.Unknown)
	writer.Write16(m.FH)
	writer.WriteU8(m.NewState)
	writer.Write16(m.Duration)
}
