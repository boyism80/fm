package dto

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

type Character struct {
	ID               uint32
	Name             string
	Gender           uint8
	SkinColor        uint8
	Face             uint32
	Hair             uint32
	Level            uint8
	Class            uint16
	Str              uint16
	Dex              uint16
	Int              uint16
	Luk              uint16
	Hp               uint16
	MaxHp            uint16
	Mp               uint16
	MaxMp            uint16
	AbilityPoint     uint16
	Exp              uint32
	FamePoint        uint16
	Map              uint32
	SpawnPoint       uint8
	Rank             uint32
	RankDiff         int32
	ClassRank        uint32
	ClassRankDiff    int32
	Mega             bool
	BaseLooks        map[int8]uint32
	Overlays         map[int8]uint32
	Weapon           uint32
	Position         types.Vector2[int16]
	Stance           uint8
	Meso             int32
	SkillPoint       uint16
	Inventory        map[constant.InventoryType]*Inventory
	Equipments       map[constant.EquipmentPartsType]*Equipment
	Skills           []*Skill
	Cooldowns        map[uint32]uint16
	QuestsStarted    []*QuestStatus
	QuestsCompleted  []*QuestStatus
	Rings            RingContainer
	RegRocks         []uint32
	Rocks            []uint32
	MonsterBookCover uint32
	QuestInfo        map[uint16]string
	MarriageId       uint32
	Random1          *stream.RandomStream
	Random2          *stream.RandomStream
	Random3          *stream.RandomStream
	BuddyCapacity    uint8
}

func (c *Character) SerializeOverview(writer *stream.StreamWriter) {
	writer.WriteU32(c.ID)
	writer.WriteStaticStr(c.Name, 13)
	writer.WriteU8(c.Gender)
	writer.WriteU8(c.SkinColor)
	writer.WriteU32(c.Face)
	writer.WriteU32(c.Hair)
	writer.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	writer.WriteU8(c.Level)
	writer.WriteU16(c.Class)
	writer.WriteU16(c.Str)
	writer.WriteU16(c.Dex)
	writer.WriteU16(c.Int)
	writer.WriteU16(c.Luk)
	writer.WriteU16(c.Hp)
	writer.WriteU16(c.MaxHp)
	writer.WriteU16(c.Mp)
	writer.WriteU16(c.MaxMp)
	writer.WriteU16(c.AbilityPoint)
	writer.WriteU16(0)
	writer.WriteU32(c.Exp)
	writer.WriteU16(c.FamePoint)
	writer.WriteU32(c.Map)
	writer.WriteU8(c.SpawnPoint)

	writer.WriteU8(c.Gender)
	writer.WriteU8(c.SkinColor)
	writer.WriteU32(c.Face)
	writer.WriteBoolean(c.Mega)
	writer.WriteU32(c.Hair)

	for parts, itemId := range c.BaseLooks {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	for parts, itemId := range c.Overlays {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	writer.WriteU32(c.Weapon)
	writer.WriteU32(0)

	isRanked := c.Rank > 0
	writer.WriteBoolean(isRanked)
	if isRanked {
		writer.WriteU32(c.Rank)
		writer.Write32(c.RankDiff)
		writer.WriteU32(c.ClassRank)
		writer.Write32(c.ClassRankDiff)
	}
}

func (c *Character) SerializeLook(writer *stream.StreamWriter) {
	writer.WriteU8(c.Gender)
	writer.WriteU8(c.SkinColor)
	writer.WriteU32(c.Face)
	writer.WriteBoolean(c.Mega)
	writer.WriteU32(c.Hair)

	for parts, itemId := range c.BaseLooks {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	for parts, itemId := range c.Overlays {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	writer.WriteU32(c.Weapon)
	writer.WriteU32(0)
}

func (c *Character) SerializeStats(writer *stream.StreamWriter) {
	writer.WriteU32(c.ID)
	writer.WriteStaticStr(c.Name, 13)
	writer.WriteU8(c.Gender)
	writer.WriteU8(c.SkinColor)
	writer.WriteU32(c.Face)
	writer.WriteU32(c.Hair)
	writer.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	writer.WriteU8(c.Level)
	writer.WriteU16(c.Class)
	writer.WriteU16(c.Str)
	writer.WriteU16(c.Dex)
	writer.WriteU16(c.Int)
	writer.WriteU16(c.Luk)
	writer.WriteU16(c.Hp)
	writer.WriteU16(c.MaxHp)
	writer.WriteU16(c.Mp)
	writer.WriteU16(c.MaxMp)
	writer.WriteU16(c.AbilityPoint)
	writer.WriteU16(c.SkillPoint)
	writer.WriteU32(c.Exp)
	writer.WriteU16(c.FamePoint)
	writer.WriteU32(c.Map)
	writer.WriteU8(c.SpawnPoint)
}

func (c *Character) SerializeInventory(writer *stream.StreamWriter) {
	writer.Write32(c.Meso)

	writer.WriteU8(c.Inventory[constant.InventoryTypeEquipment].SlotLimit)
	writer.WriteU8(c.Inventory[constant.InventoryTypeConsume].SlotLimit)
	writer.WriteU8(c.Inventory[constant.InventoryTypeInstallation].SlotLimit)
	writer.WriteU8(c.Inventory[constant.InventoryTypeETC].SlotLimit)
	writer.WriteU8(c.Inventory[constant.InventoryTypeCash].SlotLimit)

	equipParts1 := make([]constant.EquipmentPartsType, 0)
	for parts := range c.Equipments {
		if c.Equipments[parts] != nil && (parts <= 0 && parts > -100) {
			equipParts1 = append(equipParts1, parts)
		}
	}
	sort.Slice(equipParts1, func(i, j int) bool {
		return equipParts1[i] > equipParts1[j]
	})
	for _, parts := range equipParts1 {
		c.Equipments[parts].Serialize(writer, ItemSerializeOption{
			Trade:    true,
			Slot:     int16(parts),
			SlotMode: SlotEncodeActual,
		})
	}
	writer.WriteU8(0)

	equipParts2 := make([]constant.EquipmentPartsType, 0)
	for parts := range c.Equipments {
		if c.Equipments[parts] != nil && (parts <= -100 && parts > -1000) {
			equipParts2 = append(equipParts2, parts)
		}
	}
	sort.Slice(equipParts2, func(i, j int) bool {
		return equipParts2[i] > equipParts2[j]
	})
	for _, parts := range equipParts2 {
		c.Equipments[parts].Serialize(writer, ItemSerializeOption{
			Trade:    true,
			Slot:     int16(parts),
			SlotMode: SlotEncodeActual,
		})
	}
	writer.WriteU8(0)

	c.Inventory[constant.InventoryTypeEquipment].Serialize(writer)
	c.Inventory[constant.InventoryTypeConsume].Serialize(writer)
	c.Inventory[constant.InventoryTypeInstallation].Serialize(writer)
	c.Inventory[constant.InventoryTypeETC].Serialize(writer)
	c.Inventory[constant.InventoryTypeCash].Serialize(writer)
}

func (c *Character) SerializeSkills(writer *stream.StreamWriter) {
	if c.Skills == nil {
		writer.WriteU16(0)
		return
	}

	writer.WriteU16(uint16(len(c.Skills)))
	for _, skill := range c.Skills {
		writer.WriteU32(skill.ID)
		writer.WriteU32(skill.SkillLevel)

		if (skill.ID/10000)%100 > 0 && (skill.ID/10000)%10 == 2 {
			writer.WriteU32(skill.MasterLevel)
		}
	}
}

func (c *Character) SerializeCooldowns(writer *stream.StreamWriter) {
	if c.Cooldowns == nil || len(c.Cooldowns) == 0 {
		writer.WriteU16(0)
		return
	}
	keys := make([]uint32, 0, len(c.Cooldowns))
	for skillID := range c.Cooldowns {
		keys = append(keys, skillID)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	writer.WriteU16(uint16(len(keys)))
	for _, skillID := range keys {
		writer.WriteU32(skillID)
		writer.WriteU16(c.Cooldowns[skillID])
	}
}

func (c *Character) SerializeQuests(writer *stream.StreamWriter) {
	started := c.QuestsStarted
	if started == nil {
		started = []*QuestStatus{}
	}

	writer.WriteU16(uint16(len(started)))
	for _, q := range started {
		writer.WriteU16(q.QuestID)

		if len(q.MobKills) > 0 {
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

	completed := c.QuestsCompleted
	if completed == nil {
		completed = []*QuestStatus{}
	}

	writer.WriteU16(uint16(len(completed)))
	for _, q := range completed {
		writer.WriteU16(q.QuestID)
		writer.WriteDateTime(q.CompletionTime)
	}
}

func (c *Character) SerializeRings(writer *stream.StreamWriter) {
	writer.WriteU16(0)

	left := c.Rings.Left
	if left == nil {
		left = []*Ring{}
	}
	writer.WriteU16(uint16(len(left)))
	for _, ring := range left {
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteStaticStr(ring.PartnerName, 13)
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
	}

	mid := c.Rings.Mid
	if mid == nil {
		mid = []*Ring{}
	}
	writer.WriteU16(uint16(len(mid)))
	for _, ring := range mid {
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteStaticStr(ring.PartnerName, 13)
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
		writer.WriteU32(ring.ItemId)
	}

	right := c.Rings.Right
	if right == nil {
		right = []*Ring{}
	}
	writer.WriteU16(uint16(len(right)))

	for _, ring := range right {
		writer.WriteU32(c.MarriageId)
		writer.WriteU32(0)
		writer.WriteU32(0)
		writer.WriteU16(0)
		writer.WriteU32(ring.ItemId)
		writer.WriteU32(ring.ItemId)
		writer.WriteStaticStr("", 13)
		writer.WriteStaticStr("", 13)
	}
}

func (c *Character) SerializeRocks(writer *stream.StreamWriter) {
	if c.RegRocks != nil {
		for _, regRock := range c.RegRocks {
			writer.WriteU32(regRock)
		}
	}

	if c.Rocks != nil {
		for _, rock := range c.Rocks {
			writer.WriteU32(rock)
		}
	}
}

func (c *Character) SerializeMonsterBook(writer *stream.StreamWriter) {
	writer.WriteU32(c.MonsterBookCover)
	writer.WriteU8(0)
	writer.WriteU16(0)
}

func (c *Character) SerializeQuestInfo(writer *stream.StreamWriter) {
	if c.QuestInfo == nil {
		writer.WriteU16(0)
		return
	}

	writer.WriteU16(uint16(len(c.QuestInfo)))
	for questId, customData := range c.QuestInfo {
		writer.WriteU16(questId)
		if customData == "" {
			writer.WriteStr8("")
		} else {
			writer.WriteStr8(customData)
		}
	}
}

func (c *Character) Serialize(writer *stream.StreamWriter) {
	writer.WriteU64(0xFFFFFFFFFFFFFFFF)

	c.SerializeStats(writer)
	if c.BuddyCapacity == 0 {
		c.BuddyCapacity = pconst.DefaultBuddyCapacity
	}
	writer.WriteU8(c.BuddyCapacity)

	c.SerializeInventory(writer)

	c.SerializeSkills(writer)

	c.SerializeCooldowns(writer)

	c.SerializeQuests(writer)

	c.SerializeRings(writer)

	c.SerializeRocks(writer)

	c.SerializeMonsterBook(writer)

	c.SerializeQuestInfo(writer)

	writer.WriteU16(0)
}
