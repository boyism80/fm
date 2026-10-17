package game

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
)

func (ch *Character) BuildS2C() *dto.Character {
	look := ch.inventory.look()
	equipments, overlays := look.split()
	weapon := look[constant.EquipmentPartsWeapon]

	mapID := uint32(0)
	if m := ch.GetMap(); m != nil {
		mapID = m.TemplateID()
	}
	var pet uint64
	if ch.pets.active != nil {
		pet = *ch.pets.active.item.uniqueId
	}
	return &dto.Character{
		ID:            ch.GetID(),
		Pet:           pet,
		Name:          ch.name,
		Gender:        ch.look.gender,
		SkinColor:     ch.look.skinColor,
		Face:          ch.look.face,
		Hair:          ch.look.hair,
		Level:         ch.level,
		Class:         uint16(ch.class),
		Str:           ch.stats.TotalStr(),
		Dex:           ch.stats.TotalDex(),
		Int:           ch.stats.TotalInt(),
		Luk:           ch.stats.TotalLuk(),
		Hp:            uint16(ch.GetHp()),
		MaxHp:         uint16(ch.GetMaxHp()),
		Mp:            uint16(ch.GetMp()),
		MaxMp:         uint16(ch.GetMaxMp()),
		AbilityPoint:  ch.points.ap,
		Exp:           ch.exp,
		Population:    ch.stats.population,
		Map:           mapID,
		SpawnPoint:    ch.GetSpawnPoint(),
		Rank:          ch.rank.world,
		RankDiff:      ch.rank.worldDiff,
		ClassRank:     ch.rank.class,
		ClassRankDiff: ch.rank.classDiff,
		Mega:          ch.mega,
		BaseLooks:     equipments,
		Overlays:      overlays,
		Weapon:        weapon,
		Position:      ch.position,
		Stance:        ch.stance,
	}
}

func (ch *Character) BuildFullS2C() *dto.Character {
	charDTO := ch.BuildS2C()

	charDTO.SkillPoint = ch.points.sp
	if ch.wedding.marriage != nil {
		charDTO.Marriage = ch.wedding.marriage.BuildS2C()
	}
	charDTO.TeleportStones = ch.teleportStones.regular
	charDTO.VipTeleportStones = ch.teleportStones.vip
	charDTO.MonsterBookCover = ch.monsterBook.cover
	charDTO.MonsterBookCards = ch.monsterBook.cards
	if ch.quests != nil {
		charDTO.RecordExByQuest = ch.quests.RecordExWireMap()
	}
	if capacity := ch.buddies.Capacity(); capacity > 255 {
		charDTO.BuddyCapacity = 255
	} else {
		charDTO.BuddyCapacity = uint8(capacity)
	}

	charDTO.Random1 = &ch.random[0]
	charDTO.Random2 = &ch.random[1]
	charDTO.Random3 = &ch.random[2]

	charDTO.Inventory = &dto.Inventory{
		Tabs:     make(map[constant.InventoryType]*dto.InventoryTab),
		Equipped: make(map[constant.EquipmentPartsType]*dto.Equipment),
		Rings:    ch.inventory.rings.BuildS2C(),
		Meso:     ch.inventory.meso,
	}
	for invType, inv := range ch.inventory.tabs {
		if inv == nil {
			continue
		}
		invDTO := &dto.InventoryTab{
			Type:      invType,
			SlotLimit: inv.slotLimit,
		}

		invDTO.Items = make(map[int16]dto.Item)
		for slot, item := range inv.items {
			if item == nil {
				continue
			}
			invDTO.Items[slot] = item.BuildS2C()
		}
		charDTO.Inventory.Tabs[invType] = invDTO
	}

	for parts, equipment := range ch.inventory.equipped {
		if equipment == nil {
			continue
		}
		equipDTO := equipment.BuildEquipmentS2C()
		if equipDTO != nil {
			charDTO.Inventory.Equipped[parts] = equipDTO
		}
	}

	charDTO.Skills = make([]*dto.Skill, 0)
	ch.skills.ForEach(func(skillID uint32, entry *SkillEntry) {
		skillDTO := &dto.Skill{
			ID:          skillID,
			SkillLevel:  uint32(entry.Level()),
			MasterLevel: uint32(entry.masterLevel),
		}
		charDTO.Skills = append(charDTO.Skills, skillDTO)
	})

	charDTO.Cooldowns = make(map[uint32]uint16)
	ch.skills.ForEach(func(skillID uint32, entry *SkillEntry) {
		if !entry.IsCooling() {
			return
		}
		sec := int(entry.CooldownRemaining().Seconds())
		sec = min(sec, 65535)
		charDTO.Cooldowns[skillID] = uint16(sec)
	})

	charDTO.QuestsStarted = make([]*dto.QuestStatus, 0)
	charDTO.QuestsCompleted = make([]*dto.QuestStatus, 0)
	if ch.quests != nil {
		ch.quests.ForEach(func(questID uint32, qp *Quest) {
			questDTO, ok := qp.BuildS2C()
			if !ok {
				return
			}
			switch qp.status {
			case QuestStatusStarted:
				charDTO.QuestsStarted = append(charDTO.QuestsStarted, &questDTO)
			case QuestStatusCompleted:
				charDTO.QuestsCompleted = append(charDTO.QuestsCompleted, &questDTO)
			}
		})
	}

	return charDTO
}
