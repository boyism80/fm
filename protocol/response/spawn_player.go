package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type SpawnPlayer struct {
	Character         *dto.Character
	GuildName         string
	GuildLogoBG       uint16
	GuildLogoBGColor  uint8
	GuildLogo         uint16
	GuildLogoColor    uint8
	BuffStates        [4]uint32
	Diseases          [4]uint32
	SpeedBuff         uint8
	ComboCount        uint8
	WKChargeSkillId   uint32
	MorphId           uint16
	SpiritClawSkillId uint32
	ItemEffectId      uint32
	ChairId           uint32
	Balloons          uint32
	MountLevel        uint32
	MountExp          uint32
	MountFatigue      uint32
	Chalkboard        string
	Team              uint8
	CrushRings        []*dto.Ring
	FriendshipRings   []*dto.Ring
	MarriageRings     []*dto.Ring
}

func writeRings(writer *stream.StreamWriter, rings []*dto.Ring) error {
	writer.WriteU8(uint8(len(rings)))
	for _, ring := range rings {
		if ring == nil {
			continue
		}
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
		writer.WriteU64(ring.RingUniqueId)
	}
	return nil
}

func (p *SpawnPlayer) Serialize(writer *stream.StreamWriter) error {
	if p.Character == nil {
		return nil
	}

	writer.WriteU32(p.Character.ID)
	writer.WriteStr16(p.Character.Name)
	if p.GuildName == "" {
		writer.WriteU32(0)
		writer.WriteU32(0)
	} else {
		writer.WriteStr16(p.GuildName)
		writer.WriteU16(p.GuildLogoBG)
		writer.WriteU8(p.GuildLogoBGColor)
		writer.WriteU16(p.GuildLogo)
		writer.WriteU8(p.GuildLogoColor)
	}

	for i := range 4 {
		writer.WriteU32(p.BuffStates[i])
	}

	writer.WriteU16(0)
	writer.WriteU16(p.Character.Class)
	p.Character.SerializeLook(writer)
	writer.WriteU32(p.Balloons)
	writer.WriteU32(p.ItemEffectId)
	writer.WriteU32(p.ChairId)
	writer.Write16(p.Character.Position.X)
	writer.Write16(p.Character.Position.Y)
	writer.WriteU8(p.Character.Stance)
	writer.WriteU16(0)
	writer.WriteU8(0)
	writer.WriteU32(p.MountLevel)
	writer.WriteU32(p.MountExp)
	writer.WriteU32(p.MountFatigue)
	writer.WriteU8(0)
	if p.Chalkboard != "" {
		writer.WriteU8(1)
		writer.WriteStr16(p.Chalkboard)
	} else {
		writer.WriteU8(0)
	}

	if err := writeRings(writer, p.CrushRings); err != nil {
		return err
	}
	if err := writeRings(writer, p.FriendshipRings); err != nil {
		return err
	}
	if err := writeRings(writer, p.MarriageRings); err != nil {
		return err
	}

	writer.WriteU8(0)
	writer.WriteU8(p.Team)
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	return nil
}

func (s *SpawnPlayer) Opcode() uint16 {
	return 0x6E // SpawnPlayer opcode
}

func (s *SpawnPlayer) Deserialize(reader *stream.StreamReader) {
}
