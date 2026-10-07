package response

import (
	"sort"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

type CharacterInfo struct {
	Character *dto.Character
}

func (p *CharacterInfo) Serialize(writer *stream.StreamWriter) {
	writer.WriteU64(0xFFFFFFFFFFFFFFFF)

	p.serializeStats(writer)
	if p.Character.BuddyCapacity == 0 {
		p.Character.BuddyCapacity = pconst.DefaultBuddyCapacity
	}
	writer.WriteU8(p.Character.BuddyCapacity)

	p.serializeInventory(writer)

	p.serializeSkills(writer)

	p.serializeCooldowns(writer)

	p.serializeQuests(writer)

	p.serializeRings(writer)

	p.serializeRocks(writer)

	p.serializeMonsterBook(writer)

	p.serializeRecordEx(writer)

	writer.WriteU16(0)
}

func (p *CharacterInfo) serializeStats(writer *stream.StreamWriter) {
	writer.WriteU32(p.Character.ID)
	writer.WriteStaticStr(p.Character.Name, 13)
	writer.WriteU8(p.Character.Gender)
	writer.WriteU8(p.Character.SkinColor)
	writer.WriteU32(p.Character.Face)
	writer.WriteU32(p.Character.Hair)
	writer.WriteU64(p.Character.Pet)
	writer.WriteU8(p.Character.Level)
	writer.WriteU16(p.Character.Class)
	writer.WriteU16(p.Character.Str)
	writer.WriteU16(p.Character.Dex)
	writer.WriteU16(p.Character.Int)
	writer.WriteU16(p.Character.Luk)
	writer.WriteU16(p.Character.Hp)
	writer.WriteU16(p.Character.MaxHp)
	writer.WriteU16(p.Character.Mp)
	writer.WriteU16(p.Character.MaxMp)
	writer.WriteU16(p.Character.AbilityPoint)
	writer.WriteU16(p.Character.SkillPoint)
	writer.WriteU32(p.Character.Exp)
	writer.WriteU16(p.Character.Population)
	writer.WriteU32(p.Character.Map)
	writer.WriteU8(p.Character.SpawnPoint)
}

func (p *CharacterInfo) serializeInventory(writer *stream.StreamWriter) {
	writer.Write32(p.Character.Inventory.Meso)

	writer.WriteU8(p.Character.Inventory.Tabs[constant.InventoryTypeEquipment].SlotLimit)
	writer.WriteU8(p.Character.Inventory.Tabs[constant.InventoryTypeConsume].SlotLimit)
	writer.WriteU8(p.Character.Inventory.Tabs[constant.InventoryTypeInstallation].SlotLimit)
	writer.WriteU8(p.Character.Inventory.Tabs[constant.InventoryTypeETC].SlotLimit)
	writer.WriteU8(p.Character.Inventory.Tabs[constant.InventoryTypeCash].SlotLimit)

	equipParts1 := make([]constant.EquipmentPartsType, 0)
	for parts := range p.Character.Inventory.Equipped {
		if p.Character.Inventory.Equipped[parts] != nil && (parts <= 0 && parts > -100) {
			equipParts1 = append(equipParts1, parts)
		}
	}
	sort.Slice(equipParts1, func(i, j int) bool {
		return equipParts1[i] > equipParts1[j]
	})
	for _, parts := range equipParts1 {
		p.Character.Inventory.Equipped[parts].Serialize(writer, dto.ItemSerializeOption{
			Trade:    true,
			Slot:     int16(parts),
			SlotMode: dto.SlotEncodeActual,
		})
	}
	writer.WriteU8(0)

	equipParts2 := make([]constant.EquipmentPartsType, 0)
	for parts := range p.Character.Inventory.Equipped {
		if p.Character.Inventory.Equipped[parts] != nil && (parts <= -100 && parts > -1000) {
			equipParts2 = append(equipParts2, parts)
		}
	}
	sort.Slice(equipParts2, func(i, j int) bool {
		return equipParts2[i] > equipParts2[j]
	})
	for _, parts := range equipParts2 {
		p.Character.Inventory.Equipped[parts].Serialize(writer, dto.ItemSerializeOption{
			Trade:    true,
			Slot:     int16(parts),
			SlotMode: dto.SlotEncodeActual,
		})
	}
	writer.WriteU8(0)

	p.Character.Inventory.Tabs[constant.InventoryTypeEquipment].Serialize(writer)
	p.Character.Inventory.Tabs[constant.InventoryTypeConsume].Serialize(writer)
	p.Character.Inventory.Tabs[constant.InventoryTypeInstallation].Serialize(writer)
	p.Character.Inventory.Tabs[constant.InventoryTypeETC].Serialize(writer)
	p.Character.Inventory.Tabs[constant.InventoryTypeCash].Serialize(writer)
}

func (p *CharacterInfo) serializeSkills(writer *stream.StreamWriter) {
	if p.Character.Skills == nil {
		writer.WriteU16(0)
		return
	}

	writer.WriteU16(uint16(len(p.Character.Skills)))
	for _, skill := range p.Character.Skills {
		writer.WriteU32(skill.ID)
		writer.WriteU32(skill.SkillLevel)

		if (skill.ID/10000)%100 > 0 && (skill.ID/10000)%10 == 2 {
			writer.WriteU32(skill.MasterLevel)
		}
	}
}

func (p *CharacterInfo) serializeCooldowns(writer *stream.StreamWriter) {
	if p.Character.Cooldowns == nil || len(p.Character.Cooldowns) == 0 {
		writer.WriteU16(0)
		return
	}
	keys := make([]uint32, 0, len(p.Character.Cooldowns))
	for skillID := range p.Character.Cooldowns {
		keys = append(keys, skillID)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	writer.WriteU16(uint16(len(keys)))
	for _, skillID := range keys {
		writer.WriteU32(skillID)
		writer.WriteU16(p.Character.Cooldowns[skillID])
	}
}

func (p *CharacterInfo) serializeQuests(writer *stream.StreamWriter) {
	started := p.Character.QuestsStarted
	if started == nil {
		started = []*dto.QuestStatus{}
	}

	writer.WriteU16(uint16(len(started)))
	for _, q := range started {
		writer.WriteU16(q.QuestID)
		q.WriteStartedPayload(writer)
	}

	completed := p.Character.QuestsCompleted
	if completed == nil {
		completed = []*dto.QuestStatus{}
	}

	writer.WriteU16(uint16(len(completed)))
	for _, q := range completed {
		writer.WriteU16(q.QuestID)
		writer.WriteDateTime(q.CompletionTime)
	}
}

func (p *CharacterInfo) serializeRings(writer *stream.StreamWriter) {
	writer.WriteU16(0)

	left := p.Character.Inventory.Rings.Left
	if left == nil {
		left = []*dto.Ring{}
	}
	writer.WriteU16(uint16(len(left)))
	for _, ring := range left {
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteStaticStr(ring.PartnerName, 13)
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
	}

	mid := p.Character.Inventory.Rings.Mid
	if mid == nil {
		mid = []*dto.Ring{}
	}
	writer.WriteU16(uint16(len(mid)))
	for _, ring := range mid {
		writer.WriteU32(ring.PartnerChrId)
		writer.WriteStaticStr(ring.PartnerName, 13)
		writer.WriteU64(ring.RingId)
		writer.WriteU64(ring.PartnerId)
		writer.WriteU32(ring.ItemId)
	}

	right := p.Character.Inventory.Rings.Right
	if right == nil {
		right = []*dto.Ring{}
	}
	writer.WriteU16(uint16(len(right)))

	for _, ring := range right {
		writer.WriteU32(p.Character.MarriageId)
		writer.WriteU32(0)
		writer.WriteU32(0)
		writer.WriteU16(0)
		writer.WriteU32(ring.ItemId)
		writer.WriteU32(ring.ItemId)
		writer.WriteStaticStr("", 13)
		writer.WriteStaticStr("", 13)
	}
}

func (p *CharacterInfo) serializeRocks(writer *stream.StreamWriter) {
	if p.Character.RegRocks != nil {
		for _, regRock := range p.Character.RegRocks {
			writer.WriteU32(regRock)
		}
	}

	if p.Character.Rocks != nil {
		for _, rock := range p.Character.Rocks {
			writer.WriteU32(rock)
		}
	}
}

func (p *CharacterInfo) serializeMonsterBook(writer *stream.StreamWriter) {
	writer.WriteU32(p.Character.MonsterBookCover)
	writer.WriteU8(0)
	writer.WriteU16(0)
}

func (p *CharacterInfo) serializeRecordEx(writer *stream.StreamWriter) {
	if p.Character.RecordExByQuest == nil {
		writer.WriteU16(0)
		return
	}

	writer.WriteU16(uint16(len(p.Character.RecordExByQuest)))
	for questId, customData := range p.Character.RecordExByQuest {
		writer.WriteU16(questId)
		writer.WriteStr16(customData)
	}
}

func (a *CharacterInfo) Deserialize(reader *stream.StreamReader) {
	reader.Skip(8)
	if a.Character == nil {
		a.Character = &dto.Character{}
	}
	a.deserializeStats(reader)
	a.Character.BuddyCapacity = reader.ReadU8()
	a.deserializeInventory(reader)
	a.deserializeSkills(reader)
	a.deserializeCooldowns(reader)
	a.deserializeQuests(reader)
}

func (a *CharacterInfo) deserializeStats(reader *stream.StreamReader) {
	a.Character.ID = reader.ReadU32()
	a.Character.Name = reader.ReadStaticStr(13)
	a.Character.Gender = reader.ReadU8()
	a.Character.SkinColor = reader.ReadU8()
	a.Character.Face = reader.ReadU32()
	a.Character.Hair = reader.ReadU32()
	a.Character.Pet = reader.ReadU64()
	a.Character.Level = reader.ReadU8()
	a.Character.Class = reader.ReadU16()
	a.Character.Str = reader.ReadU16()
	a.Character.Dex = reader.ReadU16()
	a.Character.Int = reader.ReadU16()
	a.Character.Luk = reader.ReadU16()
	a.Character.Hp = reader.ReadU16()
	a.Character.MaxHp = reader.ReadU16()
	a.Character.Mp = reader.ReadU16()
	a.Character.MaxMp = reader.ReadU16()
	a.Character.AbilityPoint = reader.ReadU16()
	a.Character.SkillPoint = reader.ReadU16()
	a.Character.Exp = reader.ReadU32()
	a.Character.Population = reader.ReadU16()
	a.Character.Map = reader.ReadU32()
	a.Character.SpawnPoint = reader.ReadU8()
}

func (a *CharacterInfo) deserializeInventory(reader *stream.StreamReader) {
	a.Character.Inventory = &dto.Inventory{
		Meso:     reader.Read32(),
		Tabs:     make(map[constant.InventoryType]*dto.ItemContainer),
		Equipped: make(map[constant.EquipmentPartsType]*dto.Equipment),
	}
	types := []constant.InventoryType{
		constant.InventoryTypeEquipment,
		constant.InventoryTypeConsume,
		constant.InventoryTypeInstallation,
		constant.InventoryTypeETC,
		constant.InventoryTypeCash,
	}
	for _, typ := range types {
		a.Character.Inventory.Tabs[typ] = &dto.ItemContainer{Type: typ, SlotLimit: reader.ReadU8(), Items: make(map[int16]dto.Item)}
	}

	for _, offset := range []int16{0, 100} {
		for slot := reader.ReadU8(); slot != 0; slot = reader.ReadU8() {
			if item, ok := dto.NewItemFromStream(reader).(*dto.Equipment); ok {
				a.Character.Inventory.Equipped[constant.EquipmentPartsType(-int16(slot)-offset)] = item
			}
		}
	}

	for _, typ := range types {
		for slot := reader.ReadU8(); slot != 0; slot = reader.ReadU8() {
			a.Character.Inventory.Tabs[typ].Items[int16(slot)] = dto.NewItemFromStream(reader)
		}
	}
}

func (a *CharacterInfo) deserializeSkills(reader *stream.StreamReader) {
	count := reader.ReadU16()
	a.Character.Skills = make([]*dto.Skill, 0, count)
	for range count {
		skill := &dto.Skill{ID: reader.ReadU32(), SkillLevel: reader.ReadU32()}
		if (skill.ID/10000)%100 > 0 && (skill.ID/10000)%10 == 2 {
			skill.MasterLevel = reader.ReadU32()
		}
		a.Character.Skills = append(a.Character.Skills, skill)
	}
}

func (a *CharacterInfo) deserializeCooldowns(reader *stream.StreamReader) {
	count := reader.ReadU16()
	a.Character.Cooldowns = make(map[uint32]uint16, count)
	for range count {
		skillID := reader.ReadU32()
		a.Character.Cooldowns[skillID] = reader.ReadU16()
	}
}

func (a *CharacterInfo) deserializeQuests(reader *stream.StreamReader) {
	started := reader.ReadU16()
	a.Character.QuestsStarted = make([]*dto.QuestStatus, 0, started)
	for range started {
		q := &dto.QuestStatus{QuestID: reader.ReadU16(), Status: QuestWireStatusStarted}
		q.StatusRecord = reader.ReadStr16()
		a.Character.QuestsStarted = append(a.Character.QuestsStarted, q)
	}

	completed := reader.ReadU16()
	a.Character.QuestsCompleted = make([]*dto.QuestStatus, 0, completed)
	for range completed {
		q := &dto.QuestStatus{QuestID: reader.ReadU16(), Status: QuestWireStatusCompleted}
		q.CompletionTime = util.FromFileTime(reader.ReadU64())
		a.Character.QuestsCompleted = append(a.Character.QuestsCompleted, q)
	}
}
