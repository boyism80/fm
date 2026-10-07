package response

import (
	"github.com/boyism80/fm/stream"
)

type CashShopBalance struct {
	NXCash     uint32
	MaplePoint uint32
}

func (p *CashShopBalance) Opcode() uint16 {
	return 0xF9
}

func (p *CashShopBalance) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.NXCash)
	writer.WriteU32(p.MaplePoint)
	return nil
}

func (p *CashShopBalance) Deserialize(reader *stream.StreamReader) {
	p.NXCash = reader.ReadU32()
	p.MaplePoint = reader.ReadU32()
}
