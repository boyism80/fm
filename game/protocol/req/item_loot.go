package req

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type ItemLoot struct {
	Tick     uint32
	Position types.Vector2[int16]
	OID      uint32
}

func (p *ItemLoot) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *ItemLoot) Deserialize(reader *stream.StreamReader) error {
	p.Tick, _ = reader.ReadU32()
	reader.Skip(1)
	x, _ := reader.Read16()
	y, _ := reader.Read16()
	p.Position = types.Vector2[int16]{
		X: x,
		Y: y,
	}
	p.OID, _ = reader.ReadU32()
	return nil
}
