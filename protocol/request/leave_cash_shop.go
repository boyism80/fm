package request

import (
	"github.com/boyism80/fm/stream"
)

type LeaveCashShop struct{}

func (*LeaveCashShop) Opcode() byte {
	return 0x15
}

func (a *LeaveCashShop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *LeaveCashShop) Deserialize(reader *stream.StreamReader) {
}
