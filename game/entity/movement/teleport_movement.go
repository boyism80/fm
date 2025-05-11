package movement

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type TeleportMovement struct {
	Command  uint8
	Position types.Vec2[int16]
	Velocity types.Vec2[int16]
	Stance   uint8
}

func (m *TeleportMovement) GetNewState() uint8 {
	return m.Stance
}

func (m *TeleportMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Velocity.X)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Velocity.Y)
}
