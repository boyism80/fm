package request

import (
	"github.com/boyism80/fm/stream"
)

type CashShopRefresh struct{}

func (*CashShopRefresh) Opcode() byte {
	return 0xBB
}

func (a *CashShopRefresh) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (a *CashShopRefresh) Deserialize(reader *stream.StreamReader) {
}
