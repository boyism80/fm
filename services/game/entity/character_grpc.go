package entity

import (
	"time"

	"github.com/boyism80/fm/core/clock"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

func NewCharacterFromInternalProto(sender Sendable, listener CharacterListener, reply *internal.EnterGameReply, gw GameWorld) *Character {
	if listener == nil {
		panic("NewCharacterFromInternalProto: listener must not be nil")
	}
	p := reply.GetCharacter()
	var partyID, guildID *uint32
	if reply.PartyId != nil {
		v := *reply.PartyId
		partyID = &v
	}
	if reply.GuildId != nil {
		v := *reply.GuildId
		guildID = &v
	}
	ch := &Character{
		Sendable: sender,
		Listener: listener,
		LifeCore: LifeCore{
			ObjectCore: ObjectCore{
				GameWorld: gw,
				Position: types.Vector2[int16]{
					X: int16(p.GetPositionX()),
					Y: int16(p.GetPositionY()),
				},
			},
			BaseHp: p.GetMaxHp(),
			BaseMp: p.GetMaxMp(),
			Stance: uint8(p.GetStance()),
		},
		id:           p.GetCharacterId(),
		AccountID:    p.GetAccountId(),
		name:         p.GetName(),
		gender:       uint8(p.GetGender()),
		skinColor:    uint8(p.GetSkinColor()),
		face:         p.GetFace(),
		hair:         p.GetHair(),
		level:        uint8(p.GetLevel()),
		Class:        uint16(p.GetClassId()),
		Role:         constant.CharacterRole(p.GetRole()),
		hidden:       p.GetHidden(),
		BaseStats:    BaseStats{Str: uint16(p.GetStr()), Dex: uint16(p.GetDex()), Int: uint16(p.GetIntStat()), Luk: uint16(p.GetLuk())},
		AbilityPoint: uint16(p.GetAbilityPoint()),
		SkillPoint:   uint16(p.GetSkillPoint()),
		exp:          p.GetExp(),
		population:   uint16(p.GetPopulation()),
		partyID:      partyID,
		guildID:      guildID,
		buddyList:    NewBuddyList(),

		random1: stream.NewRandomStream(),
		random2: stream.NewRandomStream(),
		random3: stream.NewRandomStream(),

		regRocks: []uint32{999999999, 999999999, 999999999, 999999999, 999999999},
		rocks:    []uint32{999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999, 999999999},
	}
	ch.Buffs = NewBuffContainer(ch)
	ch.Skills = NewSkillContainer(ch)
	ch.Quests = NewQuestContainer(ch)
	ch.Summons = NewSummonContainer(ch)
	ch.Doors = NewDoorContainer(ch)
	ch.Inventory = NewInventory(ch)
	ch.Inventory.Meso = p.GetMeso()
	ch.GuildInvites = make(map[uint32]time.Time)
	ch.savedLocations = make(map[string]uint32)
	ch.keyLayout = NewKeyLayout()
	ch.LifeCore.ObjectCore.self = ch
	ch.LifeCore.ObjectCore.initTimers()
	ch.LifeCore.setHp(p.GetHp())
	ch.LifeCore.setMp(p.GetMp())

	ch.KeyLayout().LoadKeyLayoutProto(reply.GetKeyLayout())
	ch.LoadInventory(reply.GetInventory())
	ch.LoadSkills(reply.GetSkills())
	ch.LoadBuffs(reply.GetBuffs())
	ch.LoadDebuffs(reply.GetDebuffs())
	ch.LoadQuests(reply.GetQuests())
	ch.LoadSavedLocations(reply.GetSavedLocations())
	ch.BuddyList().LoadFromProto(reply.GetBuddies(), reply.GetBuddyCapacity())
	return ch
}

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
			if constant.CanEquipAt(itemID, parts) == false {
				if known := constant.EquipmentParts(itemID); len(known) > 0 && ch.Inventory.Equipped[known[0]] == nil {
					parts = known[0]
				}
			}
			if eq, ok := item.(Equipment); ok {
				ch.Inventory.Equipped[parts] = eq
			}
		} else {
			invType := constant.InventoryType(pb.GetInventoryType())
			if invType == 0 {
				invType = constant.GetInventoryTypeByItemID(itemID)
			}
			if inv, ok := ch.Inventory.Tabs[invType]; ok {
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
	if ch == nil || ch.GameWorld == nil {
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

func (ch *Character) LoadDebuffs(persisted []*internal.DebuffPersisted) {
	now := clock.Now()
	flags := constant.AllDebuffFlags()
	for _, pb := range persisted {
		var duration time.Duration
		if pb.GetEndUnixMs() != 0 {
			duration = time.UnixMilli(pb.GetEndUnixMs()).Sub(now)
			if duration <= 0 {
				continue
			}
		}
		for _, flag := range flags {
			if flag.Mask != pb.GetMask() || flag.Position != int(pb.GetPosition()) {
				continue
			}
			if ch.debuffs == nil {
				ch.debuffs = make(map[constant.DebuffFlag]*Debuff)
			}
			ch.debuffs[flag] = &Debuff{
				Flag:       flag,
				StartTime:  now,
				Duration:   duration,
				X:          int16(pb.GetX()),
				SkillID:    uint16(pb.GetSkillId()),
				SkillLevel: uint16(pb.GetSkillLevel()),
			}
		}
	}
}

func (ch *Character) LoadQuests(persisted []*internal.QuestPersisted) {
	if ch == nil {
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
		hasRecordEx := len(pb.GetRecordEx()) > 0
		if status != QuestStatusStarted && status != QuestStatusCompleted && pb.GetForfeited() == 0 && !hasRecordEx {
			continue
		}
		q := ch.Quests.Get(questID)
		if q == nil {
			q = ch.Quests.Create(questID, status)
		}
		if q == nil {
			continue
		}
		if q.Wz == nil {
			q.Wz = ch.Quests.wzDef(questID)
		}
		q.Status = status
		q.MobKills = make(map[uint32]int, len(pb.GetMobKills()))
		for mobID, kills := range pb.GetMobKills() {
			q.MobKills[mobID] = int(kills)
		}
		q.StatusRecord.WriteString(pb.GetStatusRecord())
		if ms := pb.GetDeadlineUnixMs(); ms > 0 {
			q.Deadline = time.UnixMilli(ms)
		} else {
			q.Deadline = time.Time{}
		}
		if ms := pb.GetStartTimeUnixMs(); ms > 0 {
			q.StartTime = time.UnixMilli(ms)
		} else {
			q.StartTime = time.Time{}
		}
		q.RecordEx = make(map[string]string, len(pb.GetRecordEx()))
		for key, value := range pb.GetRecordEx() {
			q.RecordEx[key] = value
		}
		if ms := pb.GetCompletionTimeUnixMs(); ms > 0 {
			q.CompletionTime = time.UnixMilli(ms)
		} else {
			q.CompletionTime = time.Time{}
		}
		q.Forfeited = int(pb.GetForfeited())
	}
}

func equipmentLooksForPersist(ch *Character) (baseLooks, overlays map[int32]uint32) {
	baseLooks = make(map[int32]uint32)
	overlays = make(map[int32]uint32)
	if ch == nil {
		return
	}
	for parts, equipment := range ch.Inventory.Equipped {
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

func (ch *Character) PersistMapID() uint32 {
	if ch == nil {
		return 0
	}
	m := ch.GetMap()
	if m == nil || m.Wz == nil {
		return 0
	}
	if m.Wz.HasForcedReturn() {
		return uint32(m.Wz.ForcedReturn)
	}
	if ch.GetHp() < 1 && m.Wz.ReturnMapId > 0 {
		return uint32(m.Wz.ReturnMapId)
	}
	return m.TemplateID()
}

func (ch *Character) ToProto(worldID uint32) *internal.CharacterSaveEntry {
	mapID := ch.PersistMapID()
	if mapID == 0 {
		return nil
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
		Meso:         ch.Inventory.Meso,
		SkillPoint:   uint32(ch.SkillPoint),
		Population:   uint32(ch.population),
	}
	return &internal.CharacterSaveEntry{
		Character:      persisted,
		BaseLooks:      baseLooks,
		Overlays:       overlays,
		Inventory:      ch.InventoryPersisted(),
		Skills:         ch.SkillsPersisted(),
		Buffs:          ch.BuffsPersisted(),
		KeyLayout:      ch.KeyLayout().ToProto(),
		Quests:         ch.QuestsPersisted(),
		SavedLocations: ch.SavedLocationsPersisted(),
	}
}

func (ch *Character) InventoryPersisted() []*internal.InventoryPersisted {
	items := make([]*internal.InventoryPersisted, 0, len(ch.Inventory.Equipped)+64)
	ownerID := ch.GetID()
	for parts, item := range ch.Inventory.Equipped {
		if item == nil {
			continue
		}
		if pb := item.ToProto(ownerID, int32(parts)); pb != nil {
			items = append(items, pb)
		}
	}
	for _, inv := range ch.Inventory.Tabs {
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

func (ch *Character) DebuffsPersisted() []*internal.DebuffPersisted {
	now := clock.Now()
	out := make([]*internal.DebuffPersisted, 0, len(ch.debuffs))
	for _, holder := range ch.debuffs {
		var endUnixMs int64
		if holder.Duration > 0 {
			end := holder.StartTime.Add(holder.Duration)
			if end.After(now) == false {
				continue
			}
			endUnixMs = end.UnixMilli()
		}
		out = append(out, &internal.DebuffPersisted{
			Mask:       holder.Flag.Mask,
			Position:   int32(holder.Flag.Position),
			X:          int32(holder.X),
			SkillId:    uint32(holder.SkillID),
			SkillLevel: uint32(holder.SkillLevel),
			EndUnixMs:  endUnixMs,
		})
	}
	return out
}

func (ch *Character) BuffsPersisted() []*internal.BuffPersisted {
	if ch == nil {
		return nil
	}
	now := clock.Now()
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

func (ch *Character) LoadSavedLocations(persisted []*internal.SavedLocationPersisted) {
	if ch == nil {
		return
	}
	ch.savedLocations = make(map[string]uint32, len(persisted))
	for _, pb := range persisted {
		if pb == nil {
			continue
		}
		name := pb.GetName()
		if name == "" {
			continue
		}
		ch.savedLocations[name] = pb.GetMapId()
	}
}

func (ch *Character) SavedLocationsPersisted() []*internal.SavedLocationPersisted {
	if ch == nil || len(ch.savedLocations) == 0 {
		return nil
	}
	out := make([]*internal.SavedLocationPersisted, 0, len(ch.savedLocations))
	for name, mapID := range ch.savedLocations {
		out = append(out, &internal.SavedLocationPersisted{
			CharacterId: ch.GetID(),
			Name:        name,
			MapId:       mapID,
		})
	}
	return out
}

func (ch *Character) QuestsPersisted() []*internal.QuestPersisted {
	if ch == nil {
		return nil
	}
	out := make([]*internal.QuestPersisted, 0)
	ch.Quests.ForEach(func(questID uint32, q *Quest) {
		if q == nil {
			return
		}
		if q.Status != QuestStatusStarted && q.Status != QuestStatusCompleted && q.Forfeited == 0 && len(q.RecordEx) == 0 {
			return
		}
		mobKills := make(map[uint32]uint32, len(q.MobKills))
		for mobID, kills := range q.MobKills {
			if kills < 0 {
				continue
			}
			mobKills[mobID] = uint32(kills)
		}
		recordEx := make(map[string]string, len(q.RecordEx))
		for key, value := range q.RecordEx {
			recordEx[key] = value
		}
		var completionTimeUnixMs int64
		if !q.CompletionTime.IsZero() {
			completionTimeUnixMs = q.CompletionTime.UnixMilli()
		}
		var deadlineUnixMs int64
		if !q.Deadline.IsZero() {
			deadlineUnixMs = q.Deadline.UnixMilli()
		}
		var startTimeUnixMs int64
		if !q.StartTime.IsZero() {
			startTimeUnixMs = q.StartTime.UnixMilli()
		}
		out = append(out, &internal.QuestPersisted{
			CharacterId:          ch.GetID(),
			QuestId:              questID,
			Status:               uint32(q.Status),
			MobKills:             mobKills,
			StatusRecord:         q.StatusRecord.AsString(),
			RecordEx:             recordEx,
			CompletionTimeUnixMs: completionTimeUnixMs,
			DeadlineUnixMs:       deadlineUnixMs,
			StartTimeUnixMs:      startTimeUnixMs,
			Forfeited:            uint32(q.Forfeited),
		})
	})
	return out
}
