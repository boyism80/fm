package entity

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
)

func (ch *Character) ToDTO() *dto.Character {
	equipments := make(map[int8]uint32)
	overlays := make(map[int8]uint32)

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
				equipments[absoluteParts] = model.GetID()
			}
		} else if absoluteParts > 100 && absoluteParts != 111 {
			adjustedParts := int8(absoluteParts - 100)
			if existingItem, exists := equipments[adjustedParts]; exists {
				overlays[adjustedParts] = existingItem
			}
			equipments[adjustedParts] = model.GetID()
		} else if _, exists := equipments[absoluteParts]; exists {
			overlays[absoluteParts] = model.GetID()
		}
	}

	var weapon uint32
	if weaponEquip := ch.Equipments[constant.EQUIPMENT_PARTS_WEAPON]; weaponEquip != nil {
		weapon = weaponEquip.Wz.GetID()
	}

	mapID := uint32(0)
	if m := ch.GetMap(); m != nil {
		mapID = m.ID
	}
	return &dto.Character{
		ID:            ch.GetID(),
		Name:          ch.name,
		Gender:        ch.gender,
		SkinColor:     ch.skinColor,
		Face:          ch.face,
		Hair:          ch.hair,
		Level:         ch.level,
		Class:         ch.Class,
		Str:           ch.GetTotalStr(),
		Dex:           ch.GetTotalDex(),
		Int:           ch.GetTotalInt(),
		Luk:           ch.GetTotalLuk(),
		Hp:            ch.Hp,
		MaxHp:         ch.GetMaxHp(),
		Mp:            ch.Mp,
		MaxMp:         ch.GetMaxMp(),
		AbilityPoint:  ch.AbilityPoint,
		Exp:           ch.exp,
		FamePoint:     ch.famePoint,
		Map:           mapID,
		SpawnPoint:    ch.spawnPoint,
		Rank:          ch.rank,
		RankDiff:      ch.rankDiff,
		ClassRank:     ch.classRank,
		ClassRankDiff: ch.classRankDiff,
		Mega:          ch.mega,
		BaseLooks:     equipments,
		Overlays:      overlays,
		Weapon:        weapon,
		Position:      ch.Position,
		Stance:        ch.Stance,
	}
}

func (ch *Character) ToFullDTO() *dto.Character {
	charDTO := ch.ToDTO()

	charDTO.Meso = ch.Meso
	charDTO.SkillPoint = ch.SkillPoint
	charDTO.MarriageId = ch.marriageId
	charDTO.RegRocks = ch.regRocks
	charDTO.Rocks = ch.rocks
	charDTO.MonsterBookCover = ch.monsterBookCover
	charDTO.QuestInfo = ch.quests
	charDTO.BuddyCapacity = 20

	charDTO.Random1 = &ch.random1
	charDTO.Random2 = &ch.random2
	charDTO.Random3 = &ch.random3

	charDTO.Inventory = make(map[constant.InventoryType]*dto.Inventory)
	for invType, inv := range ch.Inventory {
		if inv == nil {
			continue
		}
		invDTO := &dto.Inventory{
			Type:      invType,
			SlotLimit: inv.SlotLimit,
		}

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

	charDTO.Skills = make([]*dto.Skill, 0, len(ch.Skills))
	for skillID, entry := range ch.Skills {
		if entry == nil {
			continue
		}
		skillDTO := &dto.Skill{
			ID:         skillID,
			SkillLevel: uint32(entry.SkillLevel),
		}
		if entry.Skill != nil && entry.Skill.MasterLevel > 0 {
			skillDTO.MasterLevel = uint32(entry.MasterLevel)
		}
		charDTO.Skills = append(charDTO.Skills, skillDTO)
	}

	charDTO.Cooldowns = make(map[uint32]uint16)
	for skillID, entry := range ch.Skills {
		if entry == nil || !entry.IsCooling() {
			continue
		}
		sec := int(entry.CooldownRemaining().Seconds())
		sec = min(sec, 65535)
		charDTO.Cooldowns[skillID] = uint16(sec)
	}

	charDTO.QuestsStarted = make([]*dto.QuestStatus, 0)
	charDTO.QuestsCompleted = make([]*dto.QuestStatus, 0)
	for _, qs := range ch.questStatuses {
		if qs == nil {
			continue
		}
		questDTO := &dto.QuestStatus{
			QuestID:        uint16(qs.Quest.ID),
			Status:         qs.Status,
			CustomData:     qs.CustomData,
			CompletionTime: qs.CompletionTime,
		}
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

	charDTO.Rings = dto.RingContainer{
		Left:  RingsToDTO(ch.Rings.Left),
		Mid:   RingsToDTO(ch.Rings.Mid),
		Right: RingsToDTO(ch.Rings.Right),
	}

	return charDTO
}
