package request

import (
	"fmt"

	"github.com/boyism80/fm/stream"
)

type SelectCharacterWithSecondPassword struct {
	Password    string
	CharacterID uint32
}

func (*SelectCharacterWithSecondPassword) Opcode() byte { return 0x10 }

func (a *SelectCharacterWithSecondPassword) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(a.Password)
	writer.WriteU32(a.CharacterID)
	return nil
}

func (a *SelectCharacterWithSecondPassword) Deserialize(reader *stream.StreamReader) {
	a.Password = reader.ReadStr16()
	a.CharacterID = reader.ReadU32()
}

func (a *SelectCharacterWithSecondPassword) String() string {
	return fmt.Sprintf("SelectCharacterWithSecondPassword{Password:[REDACTED] CharacterID:%d}", a.CharacterID)
}
