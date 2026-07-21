package request

import (
	"github.com/boyism80/fm/stream"
)

type AllyDamage struct {
	FromOID uint32
	ToOID   uint32
}

func (*AllyDamage) Opcode() byte { return 0x99 }

func (p *AllyDamage) Serialize(_ *stream.StreamWriter) error {
	return nil
}

func (p *AllyDamage) Deserialize(reader *stream.StreamReader) {
	p.FromOID = reader.ReadU32()
	reader.Skip(4)
	p.ToOID = reader.ReadU32()
}
