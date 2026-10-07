package request

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type PetLoot struct {
	Tick     uint32
	Position types.Vector2[int16]
	OID      uint32
}

func (*PetLoot) Opcode() byte { return 0x87 }

func (p *PetLoot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU32(p.Tick)
	writer.Write16(p.Position.X)
	writer.Write16(p.Position.Y)
	writer.WriteU32(p.OID)
	return nil
}

func (p *PetLoot) Deserialize(reader *stream.StreamReader) {
	reader.Skip(1)
	p.Tick = reader.ReadU32()
	p.Position = types.Vector2[int16]{X: reader.Read16(), Y: reader.Read16()}
	p.OID = reader.ReadU32()
}
