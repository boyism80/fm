package request

import (
	"github.com/boyism80/fm/stream"
)

type LoginGame struct {
	PlayerId uint32
}

func (*LoginGame) Opcode() byte { return 0x06 }

func (a *LoginGame) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *LoginGame) Deserialize(reader *stream.StreamReader) {
	a.PlayerId = reader.ReadU32()
}
