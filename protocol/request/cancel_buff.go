package request

import (
	"github.com/boyism80/fm/stream"
)

type CancelBuff struct {
	SourceID int32
}

func (c *CancelBuff) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (c *CancelBuff) Deserialize(reader *stream.StreamReader) {
	c.SourceID = reader.Read32()
}
