package request

import (
	"github.com/boyism80/fm/stream"
)

type CashShopCoupon struct {
	Recipient string
	Code      string
	Message   string
}

func (*CashShopCoupon) Opcode() byte {
	return 0xBD
}

func (a *CashShopCoupon) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(a.Recipient)
	writer.WriteStr16(a.Code)
	if a.Recipient != "" {
		writer.WriteStr16(a.Message)
	}
	return nil
}

func (a *CashShopCoupon) Deserialize(reader *stream.StreamReader) {
	a.Recipient = reader.ReadStr16()
	a.Code = reader.ReadStr16()
	if a.Recipient != "" && reader.Remaining() >= 2 {
		a.Message = reader.ReadStr16()
	}
}
