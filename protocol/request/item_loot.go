package request

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type ItemLoot struct {
	FieldKey uint8
	Tick     uint32
	Position types.Vector2[int16]
	OID      uint32
}

func (*ItemLoot) Opcode() byte { return 0xA3 }

func (p *ItemLoot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.FieldKey)
	writer.WriteU32(p.Tick)
	writer.Write16(p.Position.X)
	writer.Write16(p.Position.Y)
	writer.WriteU32(p.OID)
	return nil
}

func (p *ItemLoot) Deserialize(reader *stream.StreamReader) {
	p.FieldKey = reader.ReadU8()
	p.Tick = reader.ReadU32()
	p.Position = types.Vector2[int16]{
		X: reader.Read16(),
		Y: reader.Read16(),
	}
	p.OID = reader.ReadU32()
}
