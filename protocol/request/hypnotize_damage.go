package request

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type HypnotizeDamage struct {
	SourceOID   uint32
	CharacterID uint32
	TargetOID   uint32
	Action      uint8
	Damage      uint32
	Left        bool
	Position    types.Vector2[int16]
}

func (*HypnotizeDamage) Opcode() byte { return 0x9B }

func (p *HypnotizeDamage) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.SourceOID)
	writer.WriteU32(p.CharacterID)
	writer.WriteU32(p.TargetOID)
	writer.WriteU8(p.Action)
	writer.WriteU32(p.Damage)
	writer.WriteBoolean(p.Left)
	writer.Write16(p.Position.X)
	writer.Write16(p.Position.Y)
	return nil
}

func (p *HypnotizeDamage) Deserialize(reader *stream.StreamReader) {
	p.SourceOID = reader.ReadU32()
	p.CharacterID = reader.ReadU32()
	p.TargetOID = reader.ReadU32()
	p.Action = reader.ReadU8()
	p.Damage = reader.ReadU32()
	p.Left = reader.ReadBool()
	p.Position.X = reader.Read16()
	p.Position.Y = reader.Read16()
}
