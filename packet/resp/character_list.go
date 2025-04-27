package resp

import (
	"github.com/boyism80/fm/dao"
	"github.com/boyism80/fm/stream"
)

type CharacterList struct {
	SecondPw   string
	Characters []dao.Character
	SlotCount  uint32
}

func writeCharacterStats(ch *dao.Character, writer *stream.StreamWriter) {
	writer.WriteU32(ch.Id)
	writer.WriteStaticStr(ch.Name, 13)
	writer.WriteU8(ch.Gender)
	writer.WriteU8(ch.SkinColor)
	writer.WriteU32(ch.Face)
	writer.WriteU32(ch.Hair)
	writer.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	writer.WriteU8(ch.Level)
	writer.WriteU16(ch.Class)
	writer.WriteU16(ch.Str)
	writer.WriteU16(ch.Dex)
	writer.WriteU16(ch.Int)
	writer.WriteU16(ch.Luk)
	writer.WriteU16(ch.Hp)
	writer.WriteU16(ch.MaxHp)
	writer.WriteU16(ch.Mp)
	writer.WriteU16(ch.MaxMp)
	writer.WriteU16(ch.AbilityPoint)

	remainingSkillPoints := ch.RemainingSkillPoints()
	if ch.IsEvan() || ch.IsResist() || ch.IsMercedes() {
		writer.WriteU8(uint8(remainingSkillPoints))

		for i, sp := range ch.SkillPoint {
			if sp == 0 {
				continue
			}
			writer.WriteU8(uint8(i))
			writer.WriteU8(uint8(sp))
		}
	} else {
		writer.WriteU16(remainingSkillPoints)
	}

	writer.WriteU32(ch.Exp)
	writer.WriteU16(ch.FamePoint)
	writer.WriteU32(ch.Map)
	writer.WriteU8(ch.SpawnPoint)
}

func writeCharacterLook(ch *dao.Character, writer *stream.StreamWriter) {
	writer.WriteU8(ch.Gender)
	writer.WriteU8(ch.SkinColor)
	writer.WriteU32(ch.Face)
	writer.WriteBoolean(ch.Mega)
	writer.WriteU32(ch.Hair)

	inventory := map[int8]uint32{
		-101: 1002577,
		-105: 1052309,
		-107: 1073255,
		-111: 1702382,
		-5:   1040010,
		-6:   1060006,
		-7:   1072005,
		-11:  1312004,
	}
	equipments := map[int8]uint32{}
	skins := map[int8]uint32{}
	for parts, itemId := range inventory {
		if parts < -127 {
			continue
		}

		realParts := int8(parts * -1)

		if realParts < 100 {
			if _, exists := equipments[realParts]; !exists {
				equipments[realParts] = itemId
			}
		} else if realParts > 100 && realParts != 111 {
			adjustedParts := int8(realParts - 100)
			if existingItem, exists := equipments[adjustedParts]; exists {
				skins[adjustedParts] = existingItem
			}
			equipments[adjustedParts] = itemId
		} else if _, exists := equipments[realParts]; exists {
			skins[realParts] = itemId
		}
	}

	for parts, itemId := range equipments {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	for parts, itemId := range skins {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	weapon := &dao.Item{
		Id: 1702382,
	}
	if weapon != nil {
		writer.WriteU32(weapon.Id)
	} else {
		writer.WriteU32(0)
	}
	writer.WriteU32(0)
}

func writeCharacter(ch *dao.Character, writer *stream.StreamWriter) {
	writeCharacterStats(ch, writer)
	writeCharacterLook(ch, writer)

	isRanked := ch.IsRanked()
	writer.WriteBoolean(isRanked)
	if isRanked {
		writer.WriteU32(ch.Rank)
		writer.Write32(ch.RankDiff)
		writer.WriteU32(ch.ClassRank)
		writer.Write32(ch.ClassRankDiff)
	}
}

func (e *CharacterList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x03)
	writer.WriteU8(0)
	writer.WriteU32(0)
	writer.WriteU8(uint8(len(e.Characters)))

	for _, ch := range e.Characters {
		writeCharacter(&ch, writer)
	}

	if e.SecondPw != "" {
		writer.WriteU8(1)
	} else {
		writer.WriteU8(2)
	}
	writer.WriteU8(0)
	writer.WriteU32(e.SlotCount)

	return nil
}

func (e *CharacterList) Deserialize(reader *stream.StreamReader) error {
	return nil
}
