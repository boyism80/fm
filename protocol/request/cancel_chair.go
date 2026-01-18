package request

import (
	"github.com/boyism80/fm/stream"
)

type CancelChair struct {
	ChairID int16
}

func (c *CancelChair) Opcode() uint16 {
	return 0x19
}

func (c *CancelChair) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (c *CancelChair) Deserialize(reader *stream.StreamReader) error {
	var err error
	if c.ChairID, err = reader.Read16(); err != nil {
		return err
	}
	return nil
}
