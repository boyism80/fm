package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type Ping struct {
}

func (a *Ping) Serialize(writer *stream.StreamWriter) error {
	err := writer.WriteU16(0x09)
	if err != nil {
		return err
	}
	return nil
}

func (a *Ping) Deserialize(reader *stream.StreamReader) error {
	return nil
}
