package entity

import (
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
)

// ToDTO converts entity Character to dto Character
func (ch *Character) ToDTO() *dto.Character {
	// Convert equipments to dto format
	equipments := make(map[int8]uint32)
	skins := make(map[int8]uint32)

	for parts, equipment := range ch.Equipments {
		if equipment == nil {
			continue
		}

		if parts < -127 {
			continue
		}

		model, ok := equipment.Wz.(*wz.Equipment)
		if !ok {
			continue
		}

		absoluteParts := int8(parts * -1)
		if absoluteParts < 100 {
			if _, exists := equipments[absoluteParts]; !exists {
				equipments[absoluteParts] = model.ID
			}
		} else if absoluteParts > 100 && absoluteParts != 111 {
			adjustedParts := int8(absoluteParts - 100)
			if existingItem, exists := equipments[adjustedParts]; exists {
				skins[adjustedParts] = existingItem
			}
			equipments[adjustedParts] = model.ID
		} else if _, exists := equipments[absoluteParts]; exists {
			skins[absoluteParts] = model.ID
		}
	}

	var weapon uint32
	if weaponEquip := ch.Equipments[constant.EQUIPMENT_PARTS_WEAPON]; weaponEquip != nil {
		weapon = weaponEquip.Wz.GetID()
	}

	return &dto.Character{
		CharacterOverview: dto.CharacterOverview{
			ID:            ch.ID,
			Name:          ch.Name,
			Gender:        ch.Gender,
			SkinColor:     ch.SkinColor,
			Face:          ch.Face,
			Hair:          ch.Hair,
			Level:         ch.Level,
			Class:         ch.Class,
			Str:           ch.Str,
			Dex:           ch.Dex,
			Int:           ch.Int,
			Luk:           ch.Luk,
			Hp:            ch.Hp,
			MaxHp:         ch.MaxHp,
			Mp:            ch.Mp,
			MaxMp:         ch.MaxMp,
			AbilityPoint:  ch.AbilityPoint,
			Exp:           ch.Exp,
			FamePoint:     ch.FamePoint,
			Map:           ch.Map,
			SpawnPoint:    ch.SpawnPoint,
			Rank:          ch.Rank,
			RankDiff:      ch.RankDiff,
			ClassRank:     ch.ClassRank,
			ClassRankDiff: ch.ClassRankDiff,
		},
		CharacterLook: dto.CharacterLook{
			Gender:     ch.Gender,
			SkinColor:  ch.SkinColor,
			Face:       ch.Face,
			Hair:       ch.Hair,
			Mega:       ch.Mega,
			Equipments: equipments,
			Skins:      skins,
			Weapon:     weapon,
		},
		Position: ch.Position,
		Stance:   ch.Stance,
	}
}

// ToFullDTO converts entity Character to full dto Character (with all data for Login packet)
func (ch *Character) ToFullDTO() *dto.Character {
	// Start with basic DTO
	charDTO := ch.ToDTO()

	// Add full data
	charDTO.Meso = ch.Meso
	charDTO.SkillPoint = ch.SkillPoint
	charDTO.MarriageId = ch.MarriageId
	charDTO.RegRocks = ch.RegRocks
	charDTO.Rocks = ch.Rocks
	charDTO.MonsterBookCover = ch.MonsterBookCover
	charDTO.QuestInfo = ch.QuestInfo
	charDTO.BuddyCapacity = 20
	charDTO.IsEvan = ch.IsEvan()
	charDTO.IsResist = ch.IsResist()
	charDTO.IsMercedes = ch.IsMercedes()

	// Convert Random streams
	charDTO.Random1 = &ch.Random1
	charDTO.Random2 = &ch.Random2
	charDTO.Random3 = &ch.Random3

	// Convert Inventory
	charDTO.Inventory = make(map[constant.InventoryType]*dto.Inventory)
	for invType, inv := range ch.Inventory {
		if inv == nil {
			continue
		}
		invDTO := &dto.Inventory{
			Type:      invType,
			SlotLimit: inv.SlotLimit,
		}

		// Convert all items (including equipment) to Items
		invDTO.Items = make(map[int16]dto.Item)
		for slot, item := range inv.Items {
			if item == nil {
				continue
			}
			itemDTO := ItemToDTO(item)
			if itemDTO != nil {
				invDTO.Items[slot] = itemDTO
			}
		}
		charDTO.Inventory[invType] = invDTO
	}

	// Convert Equipments
	charDTO.Equipments = make(map[constant.EquipmentPartsType]*dto.Equipment)
	for parts, equipment := range ch.Equipments {
		if equipment == nil {
			continue
		}
		equipDTO := equipment.ToEquipmentDTO()
		if equipDTO != nil {
			charDTO.Equipments[parts] = equipDTO
		}
	}

	// Convert Skills
	charDTO.Skills = make([]*dto.Skill, 0, len(ch.SkillsMap))
	for skill, entry := range ch.SkillsMap {
		if skill == nil || entry == nil {
			continue
		}
		skillDTO := &dto.Skill{
			ID:         skill.ID,
			SkillLevel: uint32(entry.SkillLevel),
		}
		// Check if skill needs master level
		if (skill.ID/10000)%100 > 0 && (skill.ID/10000)%10 == 2 {
			skillDTO.MasterLevel = uint32(entry.MasterLevel)
		}
		charDTO.Skills = append(charDTO.Skills, skillDTO)
	}

	// Convert Cooldowns
	charDTO.Cooldowns = make([]*dto.Cooldown, 0, len(ch.CoolDowns))
	now := time.Now().UnixMilli()
	for _, cd := range ch.CoolDowns {
		if cd == nil {
			continue
		}
		remaining := int32((cd.Length + cd.StartTime - now) / 1000)
		if remaining < 0 {
			remaining = 0
		}
		cooldownDTO := &dto.Cooldown{
			SkillId:   cd.SkillId,
			Remaining: uint16(remaining),
		}
		charDTO.Cooldowns = append(charDTO.Cooldowns, cooldownDTO)
	}

	// Convert Quests
	charDTO.QuestsStarted = make([]*dto.QuestStatus, 0)
	charDTO.QuestsCompleted = make([]*dto.QuestStatus, 0)
	for _, qs := range ch.Quests {
		if qs == nil {
			continue
		}
		questDTO := &dto.QuestStatus{
			QuestID:        uint16(qs.Quest.ID),
			Status:         qs.Status,
			CustomData:     qs.CustomData,
			CompletionTime: qs.CompletionTime,
		}
		// Convert MobKills
		if len(qs.MobKills) > 0 {
			questDTO.MobKills = make([]uint16, 0, len(qs.MobKills))
			for _, kills := range qs.MobKills {
				questDTO.MobKills = append(questDTO.MobKills, uint16(kills))
			}
		}
		if qs.Status == 1 {
			charDTO.QuestsStarted = append(charDTO.QuestsStarted, questDTO)
		} else if qs.Status == 2 {
			charDTO.QuestsCompleted = append(charDTO.QuestsCompleted, questDTO)
		}
	}

	// Convert Rings
	charDTO.Rings = dto.RingContainer{
		Left:  RingsToDTO(ch.Rings.Left),
		Mid:   RingsToDTO(ch.Rings.Mid),
		Right: RingsToDTO(ch.Rings.Right),
	}

	return charDTO
}
