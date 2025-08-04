package resp

import (
	"github.com/boyism80/fm/core/stream"
)

type Ping struct {
}

func (a *Ping) Opcode() uint16 {
	return 0x09
}

func (a *Ping) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *Ping) Deserialize(reader *stream.StreamReader) error {
	return nil
}
