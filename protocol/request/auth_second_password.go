package request

import (
	"fmt"

	"github.com/boyism80/fm/stream"
)

const (
	SecondPasswordModeRemove = 0
	SecondPasswordModeSet    = 1
)

type AuthSecondPassword struct {
	Mode     uint8
	Reserved uint32
	Password string
}

func (*AuthSecondPassword) Opcode() byte { return 0x0F }

func (a *AuthSecondPassword) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(a.Mode)
	writer.WriteU32(a.Reserved)
	writer.WriteStr16(a.Password)
	return nil
}

func (a *AuthSecondPassword) Deserialize(reader *stream.StreamReader) {
	a.Mode = reader.ReadU8()
	a.Reserved = reader.ReadU32()
	a.Password = reader.ReadStr16()
}

func (a *AuthSecondPassword) String() string {
	return fmt.Sprintf("AuthSecondPassword{Mode:%d Reserved:%d Password:[REDACTED]}", a.Mode, a.Reserved)
}
