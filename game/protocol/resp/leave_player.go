package resp

import (
	"github.com/boyism80/fm/core/stream"
)

type LeavePlayer struct {
	ID uint32
}

func (p *LeavePlayer) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.ID)

	return nil
}

func (p *LeavePlayer) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *LeavePlayer) Opcode() uint16 {
	return 0x6F
}
