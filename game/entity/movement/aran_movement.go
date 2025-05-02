package movement

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type AranMovement struct {
	Command  uint8
	NewState uint8
	Unknown  int16
	Position types.Vec2[int16]
}

func (m *AranMovement) GetNewState() uint8 {
	return m.NewState
}

func (m *AranMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Unknown)
	writer.WriteU8(m.NewState)
	writer.Write16(m.Unknown)
}
