package entity

import (
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) LoadInventory(items []*internal.InventoryPersisted) {
	for _, pb := range items {
		if pb == nil {
			continue
		}
		item, err := NewItemFromInternalProto(pb, ch.GameWorld)
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
			invType := constant.InventoryType(pb.GetInventoryType())
			if invType == 0 {
				invType = constant.GetInventoryTypeByItemID(itemID)
			}
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
		entry, err := NewSkillEntryFromInternalProto(ch, pb, ch.GameWorld)
		if err != nil {
			continue
		}
		ch.Skills.Bind(pb.GetSkillId(), entry)
	}
}

func (ch *Character) LoadBuffs(persisted []*internal.BuffPersisted) {
	if ch == nil || ch.Buffs == nil || ch.GameWorld == nil {
		return
	}
	res := ch.GameWorld.GetResources()
	if res == nil {
		return
	}
	const permanentDur = 100 * 365 * 24 * time.Hour
	for _, pb := range persisted {
		if pb == nil {
			continue
		}
		values := make(map[constant.BuffFlag]int32, len(pb.GetFlagValues()))
		for _, fv := range pb.GetFlagValues() {
			if fv == nil {
				continue
			}
			values[constant.BuffFlag{Mask: fv.GetMask(), Position: int(fv.GetPosition())}] = fv.GetValue()
		}
		if len(values) == 0 {
			continue
		}
		var dur time.Duration
		if pb.RemainingDurationMs != nil {
			if *pb.RemainingDurationMs == 0 {
				continue
			}
			dur = time.Duration(*pb.RemainingDurationMs) * time.Millisecond
		} else {
			dur = permanentDur
		}
		switch pb.GetKind() {
		case internal.BuffKind_BUFF_KIND_SKILL:
			sid := uint32(pb.GetBuffSourceId())
			wzSkill := res.GetSkill(sid)
			if wzSkill == nil {
				continue
			}
			ch.Buffs.AddBuff(wzSkill, dur, uint8(pb.GetSkillLevel()), pb.GetCauserId(), values, false)
		case internal.BuffKind_BUFF_KIND_ITEM:
			itemID := uint32(-pb.GetBuffSourceId())
			model, ok := res.Items[itemID]
			if !ok {
				continue
			}
			cw, ok := model.(*wz.Consume)
			if !ok || cw == nil {
				continue
			}
			ch.Buffs.AddItemBuff(cw, dur, values, false, false)
		default:
			continue
		}
	}
}

func (ch *Character) LoadQuests(persisted []*internal.QuestPersisted) {
	if ch == nil || ch.Quests == nil {
		return
	}
	for _, pb := range persisted {
		if pb == nil {
			continue
		}
		questID := pb.GetQuestId()
		if questID == 0 {
			continue
		}
		status := QuestStatusType(pb.GetStatus())
		if status != QuestStatusStarted && status != QuestStatusCompleted {
			continue
		}
		def := ch.GameWorld.GetResources().GetQuest(questID)
		if def == nil {
			continue
		}
		q := ch.Quests.Get(questID)
		if q == nil {
			q = ch.Quests.Create(def, status)
		}
		if q == nil {
			continue
		}
		if q.Wz == nil {
			q.Wz = def
		}
		q.Status = status
		q.MobKills = make(map[uint32]int, len(pb.GetMobKills()))
		for mobID, kills := range pb.GetMobKills() {
			q.MobKills[mobID] = int(kills)
		}
		q.StatusRecord = pb.GetStatusRecord()
		q.Unknown2 = make(map[string]string, len(pb.GetUnknown2()))
		for key, value := range pb.GetUnknown2() {
			q.Unknown2[key] = value
		}
		if ms := pb.GetCompletionTimeUnixMs(); ms > 0 {
			q.CompletionTime = time.UnixMilli(ms)
		} else {
			q.CompletionTime = time.Time{}
		}
		q.Forfeited = int(pb.GetForfeited())
		q.owner = ch
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

func (ch *Character) ToProto(worldID uint32) *internal.CharacterSaveEntry {
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
		Hidden:       ch.IsHidden(),
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
		Buffs:     ch.BuffsPersisted(),
		KeyLayout: ch.KeyLayout().ToProto(),
		Quests:    ch.QuestsPersisted(),
	}
}

func (ch *Character) InventoryPersisted() []*internal.InventoryPersisted {
	items := make([]*internal.InventoryPersisted, 0, len(ch.Equipments)+64)
	ownerID := ch.GetID()
	for parts, item := range ch.Equipments {
		if item == nil {
			continue
		}
		if pb := item.ToProto(ownerID, int32(parts)); pb != nil {
			items = append(items, pb)
		}
	}
	for _, inv := range ch.Inventory {
		items = append(items, inv.ToProto(ownerID)...)
	}
	return items
}

func (ch *Character) SkillsPersisted() []*internal.SkillPersisted {
	skills := make([]*internal.SkillPersisted, 0, 64)
	ch.Skills.ForEach(func(skillID uint32, entry *SkillEntry) {
		if entry == nil {
			return
		}
		if pb := entry.ToProto(ch.GetID(), skillID); pb != nil {
			skills = append(skills, pb)
		}
	})
	return skills
}

func (ch *Character) BuffsPersisted() []*internal.BuffPersisted {
	if ch == nil || ch.Buffs == nil {
		return nil
	}
	now := time.Now()
	out := make([]*internal.BuffPersisted, 0)
	for _, ent := range ch.Buffs.Entities() {
		if ent == nil {
			continue
		}
		flags := ent.GetFlags()
		values := ent.GetValues()
		if len(flags) == 0 || len(values) == 0 {
			continue
		}
		flagPB := make([]*internal.BuffFlagValuePersisted, 0, len(flags))
		for _, f := range flags {
			v, ok := values[f]
			if !ok {
				continue
			}
			flagPB = append(flagPB, &internal.BuffFlagValuePersisted{
				Mask:     f.Mask,
				Position: int32(f.Position),
				Value:    v,
			})
		}
		if len(flagPB) == 0 {
			continue
		}
		var kind internal.BuffKind
		var skillLevel uint32
		var causerID uint32
		var rem *uint64
		switch b := ent.(type) {
		case *SkillBuff:
			kind = internal.BuffKind_BUFF_KIND_SKILL
			skillLevel = uint32(b.SkillLevel)
			causerID = b.CauserID
			if b.BaseBuff != nil && b.Duration > 0 {
				ms := uint64(ent.RemainingDuration(now) / time.Millisecond)
				if ms == 0 {
					continue
				}
				rem = &ms
			}
		case *ItemBuff:
			kind = internal.BuffKind_BUFF_KIND_ITEM
			if b.BaseBuff != nil && b.Duration > 0 {
				ms := uint64(ent.RemainingDuration(now) / time.Millisecond)
				if ms == 0 {
					continue
				}
				rem = &ms
			}
		default:
			continue
		}
		out = append(out, &internal.BuffPersisted{
			CharacterId:         ch.GetID(),
			BuffSourceId:        ent.GetBuffID(),
			Kind:                kind,
			FlagValues:          flagPB,
			RemainingDurationMs: rem,
			SkillLevel:          skillLevel,
			CauserId:            causerID,
		})
	}
	return out
}

func (ch *Character) QuestsPersisted() []*internal.QuestPersisted {
	if ch == nil || ch.Quests == nil {
		return nil
	}
	out := make([]*internal.QuestPersisted, 0)
	ch.Quests.ForEach(func(questID uint32, q *Quest) {
		if q == nil {
			return
		}
		if q.Status != QuestStatusStarted && q.Status != QuestStatusCompleted {
			return
		}
		mobKills := make(map[uint32]uint32, len(q.MobKills))
		for mobID, kills := range q.MobKills {
			if kills < 0 {
				continue
			}
			mobKills[mobID] = uint32(kills)
		}
		unknown2 := make(map[string]string, len(q.Unknown2))
		for key, value := range q.Unknown2 {
			unknown2[key] = value
		}
		var completionTimeUnixMs int64
		if !q.CompletionTime.IsZero() {
			completionTimeUnixMs = q.CompletionTime.UnixMilli()
		}
		out = append(out, &internal.QuestPersisted{
			CharacterId:          ch.GetID(),
			QuestId:              questID,
			Status:               uint32(q.Status),
			MobKills:             mobKills,
			StatusRecord:         q.StatusRecord,
			Unknown2:             unknown2,
			CompletionTimeUnixMs: completionTimeUnixMs,
			Forfeited:            uint32(q.Forfeited),
		})
	})
	return out
}
