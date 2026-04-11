package request

import (
	"github.com/boyism80/fm/stream"
)

type CancelChair struct {
	ChairID int16
}

func (c *CancelChair) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (c *CancelChair) Deserialize(reader *stream.StreamReader) {
	c.ChairID = reader.Read16()
}
