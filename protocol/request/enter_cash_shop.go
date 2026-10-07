package request

import (
	"github.com/boyism80/fm/stream"
)

type EnterCashShop struct{}

func (*EnterCashShop) Opcode() byte {
	return 0x17
}

func (a *EnterCashShop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *EnterCashShop) Deserialize(reader *stream.StreamReader) {
}
