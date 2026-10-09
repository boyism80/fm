package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type UseEntrustedShop struct{}

func (*UseEntrustedShop) Opcode() byte { return 0x2E }

func (p *UseEntrustedShop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseEntrustedShop) Deserialize(reader *stream.StreamReader) {}

type RemoteEntrustedShop struct {
	Slot int16
}

func (*RemoteEntrustedShop) Opcode() byte { return 0x2A }

func (p *RemoteEntrustedShop) Serialize(writer *stream.StreamWriter) error {
	writer.Write16(p.Slot)
	return nil
}

func (p *RemoteEntrustedShop) Deserialize(reader *stream.StreamReader) {
	p.Slot = reader.Read16()
}

type StoreBank struct {
	Mode constant.StoreBankMode
}

func (*StoreBank) Opcode() byte { return 0x2F }

func (p *StoreBank) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	return nil
}

func (p *StoreBank) Deserialize(reader *stream.StreamReader) {
	p.Mode = constant.StoreBankMode(reader.ReadU8())
}
