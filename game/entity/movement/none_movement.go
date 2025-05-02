package movement

import "github.com/boyism80/fm/common/stream"

type NoneMovement struct {
	Command  uint8
	NewState uint8
}

func (m *NoneMovement) GetNewState() uint8 {
	return m.NewState
}

func (m *NoneMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.WriteU8(m.NewState)
}
