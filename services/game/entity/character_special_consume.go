package entity

import (
	"math/rand"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) UseSkillBook(slot int16, itemID uint32) {
	if ch.GetHp() == 0 || itemID/10000 != 228 && itemID/10000 != 229 {
		return
	}
	item := ch.Inventory.GetItem(constant.InventoryTypeConsume, slot)
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != itemID {
		return
	}
	consume, ok := item.GetModel().(*wz.Consume)
	if !ok || consume.SkillBook == nil {
		return
	}

	var target *SkillEntry
	for _, skillID := range consume.SkillBook.Skills {
		if !ch.canLearn(skillID) {
			continue
		}
		entry := ch.Skills.Learned(skillID)
		if entry == nil {
			model := ch.GameWorld.GetResources().GetSkill(skillID)
			if model == nil {
				continue
			}
			entry = NewSkillEntry(ch, model, 0, 0)
		}
		if entry.Level() < consume.SkillBook.RequiredLevel || entry.MasterLevel >= consume.SkillBook.MasterLevel {
			continue
		}
		target = entry
		break
	}
	if target == nil {
		ch.Listener.OnUseSkillBook(ch, 0, consume.SkillBook.MasterLevel, false, false)
		return
	}
	if !ch.Inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1) {
		return
	}

	success := rand.Intn(100) < consume.SkillBook.Success
	if success {
		target.SetMasterLevel(consume.SkillBook.MasterLevel)
		ch.Skills.Bind(target.Wz.ID, target)
	}
	ch.Listener.OnUseSkillBook(ch, target.Wz.ID, consume.SkillBook.MasterLevel, true, success)
}

func (ch *Character) UseSummonBag(slot int16, itemID uint32) {
	if ch.GetHp() == 0 || itemID/10000 != 210 {
		return
	}
	item := ch.Inventory.GetItem(constant.InventoryTypeConsume, slot)
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != itemID {
		return
	}
	consume, ok := item.GetModel().(*wz.Consume)
	if !ok || len(consume.SummonBag) == 0 {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil || mapInstance.TemplateID() >= 910000000 && mapInstance.TemplateID() <= 910000022 {
		return
	}

	for _, entry := range consume.SummonBag {
		if mapInstance.GameWorld.GetResources().Monsters[entry.MobID] == nil {
			return
		}
	}
	var selected []uint32
	for _, entry := range consume.SummonBag {
		if entry.Probability > 0 && rand.Intn(100) < entry.Probability {
			selected = append(selected, entry.MobID)
		}
	}
	if len(selected) == 0 || !ch.Inventory.RemoveItem(constant.InventoryTypeConsume, slot, 1) {
		return
	}

	for _, mobID := range selected {
		mapInstance.SpawnMob(mobID, ch.GetPosition(), nil, constant.MobSpawnTypeAnimate, 0, constant.CarnivalTeamNone)
	}
}
