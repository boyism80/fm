package request

import (
	"github.com/boyism80/fm/stream"
)

type SelectCharacter struct {
	CharacterId uint32
}

func (*SelectCharacter) Opcode() byte { return 0x05 }

func (a *SelectCharacter) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *SelectCharacter) Deserialize(reader *stream.StreamReader) {
	a.CharacterId = reader.ReadU32()
}
