package entity

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/data"
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

func (ch *Character) SerializeStats(writer *stream.StreamWriter) {
	writer.WriteU32(ch.ID)
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

	equipments := map[int8]uint32{}
	skins := map[int8]uint32{}
	for parts, equipment := range ch.Equipments {
		if equipment == nil {
			continue
		}

		if parts < -127 {
			continue
		}

		spec, ok := equipment.Spec.(*data.EquipmentSpec)
		if !ok {
			continue
		}

		absoluteParts := int8(parts * -1)
		if absoluteParts < 100 {
			if _, exists := equipments[absoluteParts]; !exists {
				equipments[absoluteParts] = spec.ID
			}
		} else if absoluteParts > 100 && absoluteParts != 111 {
			adjustedParts := int8(absoluteParts - 100)
			if existingItem, exists := equipments[adjustedParts]; exists {
				skins[adjustedParts] = existingItem
			}
			equipments[adjustedParts] = spec.ID
		} else if _, exists := equipments[absoluteParts]; exists {
			skins[absoluteParts] = spec.ID
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

	weapon := ch.Equipments[constant.EquipmentPartsWeapon]
	if weapon != nil {
		writer.WriteU32(weapon.Spec.GetID())
	} else {
		writer.WriteU32(0)
	}
	writer.WriteU32(0)
}

func (ch *Character) SerializeOverview(writer *stream.StreamWriter) {
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

func (ch *Character) Serialize(writer *stream.StreamWriter) {
	writer.WriteU64(0xFFFFFFFFFFFFFFFF) // flag

	// flag 0x1
	ch.SerializeStats(writer)
	writer.WriteU8(20) // buddy capacity

	// flag 0x2 ~ 0x40
	ch.SerializeInventory(writer)

	// flag 0x100
	ch.SerializeSkills(writer)

	// flag 0x8000
	ch.SerializeCooldowns(writer)

	// flag 0x200, 0x4000
	ch.SerializeQuests(writer)

	// flag 0x400, 0x800
	ch.SerializeRings(writer)

	// flag 0x1000
	ch.SerializeRocks(writer)

	// flag 0x20000, 0x10000
	ch.SerializeMonsterBook(writer)

	// flag 0x40000
	ch.SerializeQuestInfo(writer)

	writer.WriteU16(0)
}

func (ch *Character) SerializeInventory(writer *stream.StreamWriter) {
	writer.Write32(ch.Meso)

	writer.WriteU8(ch.Inventory[constant.InventoryTypeEquipment].SlotLimit)
	writer.WriteU8(ch.Inventory[constant.InventoryTypeConsume].SlotLimit)
	writer.WriteU8(ch.Inventory[constant.InventoryTypeInstallation].SlotLimit)
	writer.WriteU8(ch.Inventory[constant.InventoryTypeEtc].SlotLimit)
	writer.WriteU8(ch.Inventory[constant.InventoryTypeCash].SlotLimit)

	for parts, equipment := range ch.Equipments {
		if equipment != nil && (parts <= 0 && parts > -100) {
			equipment.Serialize(writer, true, int16(parts))
		}
	}
	writer.WriteU8(0)

	for parts, equipment := range ch.Equipments {
		if equipment != nil && (parts <= -100 && parts > -1000) {
			equipment.Serialize(writer, true, int16(parts))
		}
	}
	writer.WriteU8(0)

	ch.Inventory[constant.InventoryTypeEquipment].Serialize(writer)
	ch.Inventory[constant.InventoryTypeConsume].Serialize(writer)
	ch.Inventory[constant.InventoryTypeInstallation].Serialize(writer)
	ch.Inventory[constant.InventoryTypeEtc].Serialize(writer)
	ch.Inventory[constant.InventoryTypeCash].Serialize(writer)
}

func (ch *Character) SerializeSkills(writer *stream.StreamWriter) {
	writer.WriteU16(uint16(len(ch.SkillsMap)))

	for skill, entry := range ch.SkillsMap {
		writer.WriteU32(skill.ID)
		writer.WriteU32(uint32(entry.SkillLevel))

		if (skill.ID/10000)%100 > 0 && (skill.ID/10000)%10 == 2 {
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
		writer.WriteU16(uint16(q.Quest.ID))

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
					writer.WriteDateTime(util.GetTime(timeVal))
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
		writer.WriteU16(uint16(q.Quest.ID))
		writer.WriteDateTime(q.CompletionTime)
	}
}

func (ch *Character) SerializeRings(writer *stream.StreamWriter) {
	writer.WriteU16(0)

	writer.WriteU16(uint16(len(ch.Rings.Left)))

	for _, ring := range ch.Rings.Left {
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteStaticStr(ring.PartnerName, 13)
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
	}

	writer.WriteU16(uint16(len(ch.Rings.Mid)))

	for _, ring := range ch.Rings.Mid {
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteStaticStr(ring.PartnerName, 13)
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
		writer.WriteU32(ring.ItemId)
	}

	writer.WriteU16(uint16(len(ch.Rings.Right)))
	for _, ring := range ch.Rings.Right {
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
