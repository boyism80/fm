package request

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type DamageReactor struct {
	OID     uint32
	HitSide constant.ReactorHitSide
	Stance  int16
}

func (*DamageReactor) Opcode() byte {
	return 0xA6
}

func (p *DamageReactor) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.Write32(int32(p.HitSide))
	writer.Write16(p.Stance)
	return nil
}

func (p *DamageReactor) Deserialize(reader *stream.StreamReader) {
	p.OID = reader.ReadU32()
	p.HitSide = constant.ReactorHitSide(reader.Read32())
	p.Stance = reader.Read16()
}
