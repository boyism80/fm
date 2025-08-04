package req

import (
	"github.com/boyism80/fm/core/stream"
)

type Pong struct {
}

func (a *Pong) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *Pong) Deserialize(reader *stream.StreamReader) error {
	return nil
}
