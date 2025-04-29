package entity

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

func (c *Character) cooldowns() []*CooldownEntry {
	ret := make([]*CooldownEntry, 0, len(c.CoolDowns))
	for _, cd := range c.CoolDowns {
		if cd != nil {
			ret = append(ret, cd)
		}
	}
	return ret
}

func (c *Character) startedQuests() []*QuestStatus {
	ret := make([]*QuestStatus, 0)
	for _, qs := range c.Quests {
		if qs != nil && qs.Status == 1 { // 1 = Started
			ret = append(ret, qs)
		}
	}
	return ret
}

func (c *Character) completedQuests() []*QuestStatus {
	ret := make([]*QuestStatus, 0)
	for _, qs := range c.Quests {
		if qs != nil && qs.Status == 2 { // 2 = Completed
			ret = append(ret, qs)
		}
	}
	return ret
}

func (c *Character) getRings(onlyEquipped bool) *RingCollection {
	// 임시 더미 데이터 (실제 DB/메모리에서 로드하는건 나중에 구현)
	return &RingCollection{
		Left:  []*Ring{},
		Mid:   []*Ring{},
		Right: []*Ring{},
	}
}

func (ch *Character) SerializeStats(writer *stream.StreamWriter) {
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

func (ch *Character) SerializeLook(writer *stream.StreamWriter) {
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
	ch.SerializeStats(writer)
	ch.SerializeLook(writer)

	isRanked := ch.IsRanked()
	writer.WriteBoolean(isRanked)
	if isRanked {
		writer.WriteU32(ch.Rank)
		writer.Write32(ch.RankDiff)
		writer.WriteU32(ch.ClassRank)
		writer.Write32(ch.ClassRankDiff)
	}
}

func (ch *Character) SerializeInventory(writer *stream.StreamWriter) {
	writer.WriteU32(ch.Meso)

	writer.WriteU8(ch.Inventory[InventoryTypeEquip].SlotLimit)
	writer.WriteU8(ch.Inventory[InventoryTypeUse].SlotLimit)
	writer.WriteU8(ch.Inventory[InventoryTypeSetUp].SlotLimit)
	writer.WriteU8(ch.Inventory[InventoryTypeEtc].SlotLimit)
	writer.WriteU8(ch.Inventory[InventoryTypeCash].SlotLimit)

	equipped := ch.Inventory[InventoryTypeEquipped].NewList()
	sort.Slice(equipped, func(i, j int) bool {
		return equipped[i].Parts < equipped[j].Parts
	})

	for _, item := range equipped {
		if item.Parts < 0 && item.Parts > -100 {
			item.Serialize(writer, false, false, true, false)
		}
	}
	writer.WriteU8(0)

	for _, item := range equipped {
		if item.Parts <= -100 && item.Parts > -1000 {
			item.Serialize(writer, false, false, true, false)
		}
	}
	writer.WriteU8(0)

	for _, item := range ch.Inventory[InventoryTypeEquip].NewList() {
		item.Serialize(writer, false, false, true, false)
	}
	writer.WriteU8(0)

	for _, item := range ch.Inventory[InventoryTypeUse].NewList() {
		item.Serialize(writer, false, false, true, false)
	}
	writer.WriteU8(0)

	for _, item := range ch.Inventory[InventoryTypeSetUp].NewList() {
		item.Serialize(writer, false, false, true, false)
	}
	writer.WriteU8(0)

	for _, item := range ch.Inventory[InventoryTypeEtc].NewList() {
		if item.Parts < 100 {
			item.Serialize(writer, false, false, true, false)
		}
	}
	writer.WriteU8(0)

	for _, item := range ch.Inventory[InventoryTypeCash].NewList() {
		item.Serialize(writer, false, false, true, false)
	}
	writer.WriteU8(0)
}

func (ch *Character) SerializeSkills(writer *stream.StreamWriter) {
	writer.WriteU16(uint16(len(ch.SkillsMap)))

	for skill, entry := range ch.SkillsMap {
		writer.WriteU32(skill.Id)
		writer.WriteU32(uint32(entry.SkillLevel))

		if (skill.Id/10000)%100 > 0 && (skill.Id/10000)%10 == 2 {
			writer.WriteU32(uint32(entry.MasterLevel))
		}
	}
}

func (ch *Character) SerializeCooldowns(writer *stream.StreamWriter) {
	cd := ch.cooldowns()
	writer.WriteU16(uint16(len(cd)))

	now := time.Now().UnixMilli()

	for _, cooling := range cd {
		writer.WriteU32(cooling.SkillId)
		remaining := int32((cooling.Length + cooling.StartTime - now) / 1000)
		writer.WriteU16(uint16(remaining))
	}
}

func (ch *Character) SerializeQuests(writer *stream.StreamWriter) {
	started := ch.startedQuests()
	writer.WriteU16(uint16(len(started)))

	for _, q := range started {
		writer.WriteU16(uint16(q.Quest.Id))

		if q.HasMobKills() {
			var sb strings.Builder
			for _, kills := range q.MobKills {
				sb.WriteString(fmt.Sprintf("%03d", kills))
			}
			writer.WriteStr8(sb.String())
		} else {
			if q.CustomData != "" {
				if strings.HasPrefix(q.CustomData, "time_") {
					writer.WriteU16(9)
					writer.WriteU8(1)

					timeVal, err := strconv.ParseInt(q.CustomData[5:], 10, 64)
					if err != nil {
						timeVal = 0
					}
					writer.WriteU64(util.GetTime(timeVal))
				} else {
					writer.WriteStr8(q.CustomData)
				}
			} else {
				writer.Write([]byte{0x00, 0x00})
			}
		}
	}

	completed := ch.completedQuests()
	writer.WriteU16(uint16(len(completed)))

	for _, q := range completed {
		writer.WriteU16(uint16(q.Quest.Id))
		writer.WriteU64(util.GetTime(q.CompletionTime))
	}
}

func (ch *Character) SerializeRings(writer *stream.StreamWriter) {
	writer.WriteU16(0)

	aRing := ch.getRings(true)
	cRing := aRing.Left
	writer.WriteU16(uint16(len(cRing)))

	for _, ring := range cRing {
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteStaticStr(ring.PartnerName, 13)
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
	}

	fRing := aRing.Mid
	writer.WriteU16(uint16(len(fRing)))

	for _, ring := range fRing {
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteStaticStr(ring.PartnerName, 13)
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
		writer.WriteU32(ring.ItemId)
	}

	mRing := aRing.Right
	writer.WriteU16(uint16(len(mRing)))

	for _, ring := range mRing {
		writer.WriteU32(ch.MarriageId)

		data := GetMarriageManager().GetMarriage(ch.MarriageId)
		if data == nil {
			writer.WriteU32(0)
			writer.WriteU32(0)
			writer.WriteU16(0)
			writer.WriteU32(ring.ItemId)
			writer.WriteU32(ring.ItemId)
			writer.WriteStaticStr("", 13)
			writer.WriteStaticStr("", 13)
			continue
		}

		writer.WriteU32(data.GroomId)
		writer.WriteU32(data.BrideId)

		status := data.Status
		if status == 2 {
			status = 3
		}
		writer.WriteU16(status)

		writer.WriteU32(ring.ItemId)
		writer.WriteU32(ring.ItemId)

		writer.WriteStaticStr(data.GroomName, 13)
		writer.WriteStaticStr(data.BrideName, 13)
	}
}

func (ch *Character) SerializeRocks(writer *stream.StreamWriter) {
	for _, regRock := range ch.RegRocks {
		writer.WriteU32(regRock)
	}

	for _, rock := range ch.Rocks {
		writer.WriteU32(rock)
	}
}

func (ch *Character) SerializeMonsterBook(writer *stream.StreamWriter) {
	writer.WriteU32(ch.MonsterBookCover)
	writer.WriteU8(0)

	if ch.MonsterBook != nil {
		ch.MonsterBook.Serialize(writer)
	} else {
		writer.WriteU16(0)
	}
}

func (ch *Character) SerializeQuestInfo(writer *stream.StreamWriter) {
	writer.WriteU16(uint16(len(ch.QuestInfo)))

	for questId, customData := range ch.QuestInfo {
		writer.WriteU16(questId)
		if customData == "" {
			writer.WriteStr8("")
		} else {
			writer.WriteStr8(customData)
		}
	}
}
