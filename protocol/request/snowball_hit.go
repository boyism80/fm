package request

import "github.com/boyism80/fm/stream"

// SnowballHit reports an attack on a snowball (0 bottom, 1 top) or a snowman (2 bottom, 3 top).
type SnowballHit struct {
	Target uint8
	Damage int16
	Delay  int16
}

func (*SnowballHit) Opcode() byte {
	return 0xAA
}

func (p *SnowballHit) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.Target)
	writer.Write16(p.Damage)
	writer.Write16(p.Delay)
	return nil
}

func (p *SnowballHit) Deserialize(reader *stream.StreamReader) {
	p.Target = reader.ReadU8()
	p.Damage = reader.Read16()
	p.Delay = reader.Read16()
}
