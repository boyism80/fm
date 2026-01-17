package response

import (
	"github.com/boyism80/fm/stream"
)

type ConfirmShopTransaction struct {
	Code byte
}

func (p *ConfirmShopTransaction) Opcode() uint16 {
	return 0xE7
}

func (p *ConfirmShopTransaction) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.Code)
	return nil
}

func (p *ConfirmShopTransaction) Deserialize(reader *stream.StreamReader) error {
	return nil
}
