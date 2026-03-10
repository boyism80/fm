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

func (c *CancelBuff) Deserialize(reader *stream.StreamReader) error {
	var err error
	if c.SourceID, err = reader.Read32(); err != nil {
		return err
	}
	return nil
}
