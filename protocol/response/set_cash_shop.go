package response

import (
	"github.com/boyism80/fm/stream"
)

type SetCashShop struct {
	CharacterInfo
	AccountName string
}

func (p *SetCashShop) Opcode() uint16 {
	return 0x56
}

func (p *SetCashShop) Serialize(writer *stream.StreamWriter) error {
	p.CharacterInfo.Serialize(writer)
	writer.WriteStr16(p.AccountName)
	writer.WriteU16(0)
	writer.WriteU8(0)
	writer.Write(make([]byte, 1080))
	writer.WriteU16(0)
	writer.WriteU16(0)
	writer.WriteU8(0)
	return nil
}

func (p *SetCashShop) Deserialize(reader *stream.StreamReader) {
	p.CharacterInfo.Deserialize(reader)
}
