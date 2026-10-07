package response

import (
	"github.com/boyism80/fm/stream"
)

type CharacterProfilePet struct {
	ItemID      uint32
	Name        string
	Level       uint8
	Closeness   uint16
	Fullness    uint8
	Skills      uint16
	EquipItemID uint32
}

type CharacterProfileMount struct {
	Level   uint32
	Exp     uint32
	Fatigue uint32
}

type CharacterProfile struct {
	CharacterID      uint32
	Level            uint8
	Job              uint16
	Fame             uint16
	Married          bool
	GuildName        string
	AllianceName     string
	Self             bool
	Pet              *CharacterProfilePet
	Mount            *CharacterProfileMount
	Wishlist         []uint32
	BookLevel        uint32
	BookNormalCards  uint32
	BookSpecialCards uint32
	BookCoverMobID   uint32
}

func (p *CharacterProfile) Opcode() uint16 {
	return 0x2C
}

func (p *CharacterProfile) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(p.Level)
	writer.WriteU16(p.Job)
	writer.WriteU16(p.Fame)
	writer.WriteBoolean(p.Married)
	writer.WriteStr16(p.GuildName)
	writer.WriteStr16(p.AllianceName)
	writer.WriteBoolean(p.Self)
	writer.WriteBoolean(p.Pet != nil)
	if p.Pet != nil {
		writer.WriteU32(p.Pet.ItemID)
		writer.WriteStr16(p.Pet.Name)
		writer.WriteU8(p.Pet.Level)
		writer.WriteU16(p.Pet.Closeness)
		writer.WriteU8(p.Pet.Fullness)
		writer.WriteU16(p.Pet.Skills)
		writer.WriteU32(p.Pet.EquipItemID)
	}
	writer.WriteBoolean(p.Mount != nil)
	if p.Mount != nil {
		writer.WriteU32(p.Mount.Level)
		writer.WriteU32(p.Mount.Exp)
		writer.WriteU32(p.Mount.Fatigue)
	}
	writer.WriteU8(uint8(len(p.Wishlist)))
	for _, sn := range p.Wishlist {
		writer.WriteU32(sn)
	}
	writer.WriteU32(p.BookLevel)
	writer.WriteU32(p.BookNormalCards)
	writer.WriteU32(p.BookSpecialCards)
	writer.WriteU32(p.BookNormalCards + p.BookSpecialCards)
	writer.WriteU32(p.BookCoverMobID)
	return nil
}

func (p *CharacterProfile) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.Level = reader.ReadU8()
	p.Job = reader.ReadU16()
	p.Fame = reader.ReadU16()
	p.Married = reader.ReadBool()
	p.GuildName = reader.ReadStr16()
	p.AllianceName = reader.ReadStr16()
	p.Self = reader.ReadBool()
	if reader.ReadBool() {
		p.Pet = &CharacterProfilePet{
			ItemID:      reader.ReadU32(),
			Name:        reader.ReadStr16(),
			Level:       reader.ReadU8(),
			Closeness:   reader.ReadU16(),
			Fullness:    reader.ReadU8(),
			Skills:      reader.ReadU16(),
			EquipItemID: reader.ReadU32(),
		}
	}
	if reader.ReadBool() {
		p.Mount = &CharacterProfileMount{
			Level:   reader.ReadU32(),
			Exp:     reader.ReadU32(),
			Fatigue: reader.ReadU32(),
		}
	}
	count := int(reader.ReadU8())
	p.Wishlist = make([]uint32, 0, count)
	for i := 0; i < count; i++ {
		p.Wishlist = append(p.Wishlist, reader.ReadU32())
	}
	p.BookLevel = reader.ReadU32()
	p.BookNormalCards = reader.ReadU32()
	p.BookSpecialCards = reader.ReadU32()
	reader.ReadU32()
	p.BookCoverMobID = reader.ReadU32()
}
