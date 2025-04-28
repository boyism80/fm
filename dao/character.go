package dao

import "github.com/boyism80/fm/stream"

type Character struct {
	Life
	Id            uint32
	Name          string
	Gender        uint8
	SkinColor     uint8
	Face          uint32
	Hair          uint32
	Level         uint8
	Class         uint16
	Rank          uint32
	RankDiff      int32
	ClassRank     uint32
	ClassRankDiff int32
	Admin         bool
	Str           uint16
	Dex           uint16
	Int           uint16
	Luk           uint16
	Hp            uint16
	MaxHp         uint16
	Mp            uint16
	MaxMp         uint16
	AbilityPoint  uint16
	SkillPoint    []uint16
	Exp           uint32
	FamePoint     uint16
	Map           uint32
	SpawnPoint    uint8
	Mega          bool
}

func (ch *Character) IsRanked() bool {
	if ch.Admin {
		return false
	}

	if ch.Level < 30 {
		return false
	}

	return true
}

func (ch *Character) IsEvan() bool {
	return ch.Class == 2001 || (ch.Class >= 2200 && ch.Class <= 2218)
}

func (ch *Character) IsKOC() bool {
	return ch.Class >= 1000 && ch.Class < 2000
}

func (ch *Character) IsMercedes() bool {
	return ch.Class == 2002 || (ch.Class >= 2300 && ch.Class <= 2312)
}

func (ch *Character) IsDemon() bool {
	return ch.Class == 3001 || (ch.Class >= 3100 && ch.Class <= 3112)
}

func (ch *Character) IsAran() bool {
	return (ch.Class >= 2000 && ch.Class <= 2112) && ch.Class != 2001 && ch.Class != 2002
}

func (ch *Character) IsResist() bool {
	return ch.Class >= 3000 && ch.Class <= 3512
}

func (ch *Character) IsAdventurer() bool {
	return ch.Class < 1000
}

func (ch *Character) IsCannon() bool {
	return ch.Class == 1 || ch.Class == 501 || (ch.Class >= 530 && ch.Class <= 532)
}

func (ch *Character) RemainingSkillPoints() uint16 {
	ret := 0
	for _, sp := range ch.SkillPoint {
		if sp > 0 {
			ret++
		}
	}
	return uint16(ret)
}

func (ch *Character) serializeStats(writer *stream.StreamWriter) {
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

func (ch *Character) serializeLook(writer *stream.StreamWriter) {
	writer.WriteU8(ch.Gender)
	writer.WriteU8(ch.SkinColor)
	writer.WriteU32(ch.Face)
	writer.WriteBoolean(ch.Mega)
	writer.WriteU32(ch.Hair)

	inventory := map[int8]uint32{
		// -101: 1002577,
		// -105: 1052309,
		// -107: 1073255,
		// -111: 1702382,
		// -5:   1040010,
		// -6:   1060006,
		// -7:   1072005,
		// -11:  1312004,
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

	weapon := &Item{
		Id: 1702382,
	}
	if weapon != nil {
		writer.WriteU32(weapon.Id)
	} else {
		writer.WriteU32(0)
	}
	writer.WriteU32(0)
}

func (ch *Character) Serialize(writer *stream.StreamWriter) {
	ch.serializeStats(writer)
	ch.serializeLook(writer)

	isRanked := ch.IsRanked()
	writer.WriteBoolean(isRanked)
	if isRanked {
		writer.WriteU32(ch.Rank)
		writer.Write32(ch.RankDiff)
		writer.WriteU32(ch.ClassRank)
		writer.Write32(ch.ClassRankDiff)
	}
}
