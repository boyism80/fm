package request

import (
	"github.com/boyism80/fm/stream"
)

type Pong struct {
}

func (a *Pong) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *Pong) Deserialize(reader *stream.StreamReader) {
}
