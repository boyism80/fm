package response

import (
	"github.com/boyism80/fm/stream"
)

type CancelChair struct {
	ChairID int16
}

func (c *CancelChair) Opcode() uint16 {
	return 0x96
}

func (c *CancelChair) Serialize(writer *stream.StreamWriter) error {
	if c.ChairID == -1 {
		writer.WriteU8(0)
	} else {
		writer.WriteU8(1)
		writer.Write16(c.ChairID)
	}
	return nil
}

func (c *CancelChair) Deserialize(reader *stream.StreamReader) error {
	return nil
}
