package response

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/stream"
)

// DamageCharacter shows the others on the map a character taking a hit.
type DamageCharacter struct {
	CharacterID uint32
	Type        constant.IncomingHitType
	Damage      int32
	MobID       uint32
	Direction   uint8
	// FakeSkillID is the skill a blocked hit (-1) shows, 0 for none.
	FakeSkillID uint32
}

func (p *DamageCharacter) Opcode() uint16 { return 0x89 }

func (p *DamageCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.Write8(int8(p.Type))
	writer.Write32(p.Damage)
	writer.WriteU32(p.MobID)
	writer.WriteU8(p.Direction)
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.Write32(p.Damage)
	if p.FakeSkillID > 0 {
		writer.WriteU32(p.FakeSkillID)
	}
	return nil
}

func (p *DamageCharacter) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.Type = constant.IncomingHitType(reader.Read8())
	p.Damage = reader.Read32()
	p.MobID = reader.ReadU32()
	p.Direction = reader.ReadU8()
	reader.ReadU8()
	reader.ReadU8()
	reader.Read32()
	if reader.Remaining() >= 4 {
		p.FakeSkillID = reader.ReadU32()
	}
}
