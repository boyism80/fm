package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type ActivePet struct {
	ItemID   uint32
	Name     string
	SN       uint64
	Position types.Vector2[int16]
	Stance   uint8
	Foothold int16
}

func (p *ActivePet) Serialize(writer *stream.StreamWriter) {
	writer.WriteU32(p.ItemID)
	writer.WriteStr16(p.Name)
	writer.WriteU64(p.SN)
	writer.Write16(p.Position.X)
	writer.Write16(p.Position.Y)
	writer.WriteU8(p.Stance)
	writer.Write16(p.Foothold)
}

func (p *ActivePet) Deserialize(reader *stream.StreamReader) {
	p.ItemID = reader.ReadU32()
	p.Name = reader.ReadStr16()
	p.SN = reader.ReadU64()
	p.Position.X = reader.Read16()
	p.Position.Y = reader.Read16()
	p.Stance = reader.ReadU8()
	p.Foothold = reader.Read16()
}
