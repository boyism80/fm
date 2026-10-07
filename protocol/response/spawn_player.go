package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
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
	HasTeam           bool
	Team              constant.CarnivalTeam
	CrushRing         *dto.Ring
	FriendshipRing    *dto.Ring
	MarriageRing      *dto.MarriageRing
	Pet               *dto.ActivePet
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
	for i := len(p.BuffStates) - 1; i >= 0; i-- {
		buff := p.BuffStates[i]
		if i == 3 {
			if (buff & constant.BuffFlagSpeed.Mask) != 0 {
				writer.WriteU8(p.SpeedBuff)
			}
			if (buff & constant.BuffFlagCombo.Mask) != 0 {
				writer.WriteU8(p.ComboCount)
			}
			if (buff & constant.BuffFlagWkCharge.Mask) != 0 {
				writer.WriteU32(p.WKChargeSkillId)
			}
		} else if i == 2 {
			if (buff & constant.BuffFlagMorph.Mask) != 0 {
				writer.WriteU16(p.MorphId)
			}
			if (buff & constant.BuffFlagSpiritClaw.Mask) != 0 {
				writer.WriteU32(p.SpiritClawSkillId)
			}
		}
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
	writer.WriteBoolean(p.Pet != nil)
	if p.Pet != nil {
		p.Pet.Serialize(writer)
	}
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

	writeLookRing(writer, p.CrushRing)
	writeLookRing(writer, p.FriendshipRing)
	writeMarriageRing(writer, p.MarriageRing)

	writer.WriteU8(0)
	if p.HasTeam {
		writer.Write8(int8(p.Team))
	}
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	return nil
}

func (s *SpawnPlayer) Opcode() uint16 {
	return 0x6E
}

func (s *SpawnPlayer) Deserialize(reader *stream.StreamReader) {
}
