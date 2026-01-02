package dto

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

// Character represents full character data for protocol
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
	Cooldowns        []*Cooldown
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

// SerializeOverview serializes character overview data
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
	writer.WriteU16(0) // Skill points (simplified)
	writer.WriteU32(c.Exp)
	writer.WriteU16(c.FamePoint)
	writer.WriteU32(c.Map)
	writer.WriteU8(c.SpawnPoint)

	// Look
	writer.WriteU8(c.Gender)
	writer.WriteU8(c.SkinColor)
	writer.WriteU32(c.Face)
	writer.WriteBoolean(c.Mega)
	writer.WriteU32(c.Hair)

	// Equipments
	for parts, itemId := range c.BaseLooks {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	// Overlays
	for parts, itemId := range c.Overlays {
		writer.Write8(parts)
		writer.WriteU32(itemId)
	}
	writer.WriteU8(0xFF)

	writer.WriteU32(c.Weapon)
	writer.WriteU32(0)

	// Rank info
	isRanked := c.Rank > 0
	writer.WriteBoolean(isRanked)
	if isRanked {
		writer.WriteU32(c.Rank)
		writer.Write32(c.RankDiff)
		writer.WriteU32(c.ClassRank)
		writer.Write32(c.ClassRankDiff)
	}
}

// SerializeLook serializes character look data
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

// SerializeStats serializes character stats
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

// SerializeInventory serializes character inventory
func (c *Character) SerializeInventory(writer *stream.StreamWriter) {
	writer.Write32(c.Meso)

	// Entity directly accesses inventory without nil check
	writer.WriteU8(c.Inventory[constant.INVENTORY_TYPE_EQUIPMENT].SlotLimit)
	writer.WriteU8(c.Inventory[constant.INVENTORY_TYPE_CONSUME].SlotLimit)
	writer.WriteU8(c.Inventory[constant.INVENTORY_TYPE_INSTALLATION].SlotLimit)
	writer.WriteU8(c.Inventory[constant.INVENTORY_TYPE_ETC].SlotLimit)
	writer.WriteU8(c.Inventory[constant.INVENTORY_TYPE_CASH].SlotLimit)

	// Serialize equipments (parts <= 0 && parts > -100) - sorted by parts
	equipParts1 := make([]constant.EquipmentPartsType, 0)
	for parts := range c.Equipments {
		if c.Equipments[parts] != nil && (parts <= 0 && parts > -100) {
			equipParts1 = append(equipParts1, parts)
		}
	}
	sort.Slice(equipParts1, func(i, j int) bool {
		return equipParts1[i] > equipParts1[j] // descending order (most negative first)
	})
	for _, parts := range equipParts1 {
		c.Equipments[parts].Serialize(writer, true, int16(parts))
	}
	writer.WriteU8(0)

	// Serialize equipments (parts <= -100 && parts > -1000) - sorted by parts
	equipParts2 := make([]constant.EquipmentPartsType, 0)
	for parts := range c.Equipments {
		if c.Equipments[parts] != nil && (parts <= -100 && parts > -1000) {
			equipParts2 = append(equipParts2, parts)
		}
	}
	sort.Slice(equipParts2, func(i, j int) bool {
		return equipParts2[i] > equipParts2[j] // descending order (most negative first)
	})
	for _, parts := range equipParts2 {
		c.Equipments[parts].Serialize(writer, true, int16(parts))
	}
	writer.WriteU8(0)

	// Serialize inventories (entity directly calls Serialize without nil check)
	c.Inventory[constant.INVENTORY_TYPE_EQUIPMENT].Serialize(writer)
	c.Inventory[constant.INVENTORY_TYPE_CONSUME].Serialize(writer)
	c.Inventory[constant.INVENTORY_TYPE_INSTALLATION].Serialize(writer)
	c.Inventory[constant.INVENTORY_TYPE_ETC].Serialize(writer)
	c.Inventory[constant.INVENTORY_TYPE_CASH].Serialize(writer)
}

// SerializeSkills serializes character skills
func (c *Character) SerializeSkills(writer *stream.StreamWriter) {
	if c.Skills == nil {
		writer.WriteU16(0)
		return
	}

	writer.WriteU16(uint16(len(c.Skills)))
	for _, skill := range c.Skills {
		writer.WriteU32(skill.ID)
		writer.WriteU32(skill.SkillLevel)

		// Check if skill needs master level
		if (skill.ID/10000)%100 > 0 && (skill.ID/10000)%10 == 2 {
			writer.WriteU32(skill.MasterLevel)
		}
	}
}

// SerializeCooldowns serializes character cooldowns
func (c *Character) SerializeCooldowns(writer *stream.StreamWriter) {
	if c.Cooldowns == nil {
		writer.WriteU16(0)
		return
	}

	writer.WriteU16(uint16(len(c.Cooldowns)))
	for _, cd := range c.Cooldowns {
		writer.WriteU32(cd.SkillId)
		writer.WriteU16(cd.Remaining)
	}
}

// SerializeQuests serializes character quests
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

// SerializeRings serializes character rings
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
	// Note: Right rings serialization requires marriage data which is not in DTO
	// This is a simplified version
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

// SerializeRocks serializes character rocks
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

// SerializeMonsterBook serializes character monster book
func (c *Character) SerializeMonsterBook(writer *stream.StreamWriter) {
	writer.WriteU32(c.MonsterBookCover)
	writer.WriteU8(0)
	writer.WriteU16(0) // Simplified - monster book data not in DTO
}

// SerializeQuestInfo serializes character quest info
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

// Serialize serializes full character data (for Login packet)
func (c *Character) Serialize(writer *stream.StreamWriter) {
	writer.WriteU64(0xFFFFFFFFFFFFFFFF) // flag

	// flag 0x1
	c.SerializeStats(writer)
	if c.BuddyCapacity == 0 {
		c.BuddyCapacity = 20
	}
	writer.WriteU8(c.BuddyCapacity)

	// flag 0x2 ~ 0x40
	c.SerializeInventory(writer)

	// flag 0x100
	c.SerializeSkills(writer)

	// flag 0x8000
	c.SerializeCooldowns(writer)

	// flag 0x200, 0x4000
	c.SerializeQuests(writer)

	// flag 0x400, 0x800
	c.SerializeRings(writer)

	// flag 0x1000
	c.SerializeRocks(writer)

	// flag 0x20000, 0x10000
	c.SerializeMonsterBook(writer)

	// flag 0x40000
	c.SerializeQuestInfo(writer)

	writer.WriteU16(0)
}
