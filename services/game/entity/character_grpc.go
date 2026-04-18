package entity

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) LoadInventory(items []*internal.InventoryPersisted) {
	for _, pb := range items {
		if pb == nil {
			continue
		}
		item, err := NewItemFromInternalProto(pb, ch.Context)
		if err != nil {
			continue
		}
		slot := int16(pb.GetSlot())
		itemID := pb.GetItemId()
		if slot < 0 {
			parts := constant.EquipmentPartsType(slot)
			if eq, ok := item.(Equipment); ok {
				ch.Equipments[parts] = eq
			}
		} else {
			invType := constant.GetInventoryTypeByItemID(itemID)
			if inv, ok := ch.Inventory[invType]; ok {
				inv.Items[slot] = item
			}
		}
	}
}

func (ch *Character) LoadSkills(skills []*internal.SkillPersisted) {
	for _, pb := range skills {
		if pb == nil {
			continue
		}
		entry, err := NewSkillEntryFromInternalProto(ch, pb, ch.Context)
		if err != nil {
			continue
		}
		ch.Skills.Bind(pb.GetSkillId(), entry)
	}
}

func equipmentLooksForPersist(ch *Character) (baseLooks, overlays map[int32]uint32) {
	baseLooks = make(map[int32]uint32)
	overlays = make(map[int32]uint32)
	if ch == nil {
		return
	}
	for parts, equipment := range ch.Equipments {
		if equipment == nil || parts < -127 {
			continue
		}
		model, ok := equipment.GetModel().(wz.Equipment)
		if !ok {
			continue
		}
		absoluteParts := int32(parts * -1)
		if absoluteParts < 100 {
			if _, exists := baseLooks[absoluteParts]; !exists {
				baseLooks[absoluteParts] = model.GetID()
			}
		} else if absoluteParts > 100 && absoluteParts != 111 {
			adjustedParts := absoluteParts - 100
			if existing, exists := baseLooks[adjustedParts]; exists {
				overlays[adjustedParts] = existing
			}
			baseLooks[adjustedParts] = model.GetID()
		} else if _, exists := baseLooks[absoluteParts]; exists {
			overlays[absoluteParts] = model.GetID()
		}
	}
	return
}

func (ch *Character) ToPersisted(worldID uint32) *internal.CharacterSaveEntry {
	if ch == nil {
		return nil
	}
	mapID := uint32(0)
	if m := ch.GetMap(); m != nil {
		mapID = m.GetMapID()
	}
	baseLooks, overlays := equipmentLooksForPersist(ch)
	persisted := &internal.CharacterPersisted{
		CharacterId:  ch.GetID(),
		AccountId:    ch.AccountID,
		WorldId:      worldID,
		Name:         ch.GetName(),
		Gender:       uint32(ch.GetGender()),
		SkinColor:    uint32(ch.GetSkinColor()),
		Face:         ch.GetFace(),
		Hair:         ch.GetHair(),
		Level:        uint32(ch.GetLevel()),
		ClassId:      uint32(ch.Class),
		Role:         uint32(ch.Role),
		Str:          uint32(ch.BaseStats.Str),
		Dex:          uint32(ch.BaseStats.Dex),
		IntStat:      uint32(ch.BaseStats.Int),
		Luk:          uint32(ch.BaseStats.Luk),
		Hp:           ch.GetHp(),
		MaxHp:        ch.BaseHp,
		Mp:           ch.GetMp(),
		MaxMp:        ch.BaseMp,
		AbilityPoint: uint32(ch.AbilityPoint),
		Exp:          ch.GetExp(),
		MapId:        mapID,
		SpawnPoint:   uint32(ch.GetSpawnPoint()),
		PositionX:    int32(ch.Position.X),
		PositionY:    int32(ch.Position.Y),
		Stance:       uint32(ch.Stance),
		Meso:         ch.Meso,
		SkillPoint:   uint32(ch.SkillPoint),
	}
	return &internal.CharacterSaveEntry{
		Character: persisted,
		BaseLooks: baseLooks,
		Overlays:  overlays,
		Inventory: ch.InventoryPersisted(),
		Skills:    ch.SkillsPersisted(),
		KeyLayout: ch.KeyLayout().ToPersisted(),
	}
}

func (ch *Character) InventoryPersisted() []*internal.InventoryPersisted {
	items := make([]*internal.InventoryPersisted, 0, len(ch.Equipments)+64)
	ownerID := ch.GetID()
	for parts, item := range ch.Equipments {
		if item == nil {
			continue
		}
		if pb := item.ToPersisted(ownerID, int32(parts)); pb != nil {
			items = append(items, pb)
		}
	}
	for _, inv := range ch.Inventory {
		items = append(items, inv.ToPersisted(ownerID)...)
	}
	return items
}

func (ch *Character) SkillsPersisted() []*internal.SkillPersisted {
	skills := make([]*internal.SkillPersisted, 0, 64)
	ch.Skills.ForEach(func(skillID uint32, entry *SkillEntry) {
		if entry == nil {
			return
		}
		if pb := entry.ToPersisted(ch.GetID(), skillID); pb != nil {
			skills = append(skills, pb)
		}
	})
	return skills
}
