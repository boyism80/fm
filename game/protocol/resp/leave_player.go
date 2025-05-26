package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type LeavePlayer struct {
	ID uint32
}

func (p *LeavePlayer) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x6F)
	writer.WriteU32(p.ID)

	return nil
}

func (p *LeavePlayer) Deserialize(reader *stream.StreamReader) error {
	return nil
}
