package game

import (
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/response"
)

const (
	inventorySlotDefault = 32
	inventorySlotMax     = 96
)

// ChangeClass sets the class; an advancement also brings the class change hook's gains and shows the others the effect.
func (ch *Character) ChangeClass(newClass constant.ClassType, advance bool) {
	oldClass := ch.class
	ch.class = newClass

	ch.grantClassChangePoints(newClass)

	ch.listener.OnClassChange(ch, oldClass, newClass, ch.points.sp, ch.points.ap)
	ch.NotifyLevelClassChanged()

	if !advance {
		return
	}
	ch.openFourthClassSkills(newClass)
	if m := ch.GetMap(); m != nil {
		if root := m.GetLuaRoot(); root != nil {
			if thread, err := luax.NewThread(root, constant.CharacterHookScriptPath); err == nil {
				luax.CallAsync(nil, root, thread, "on_class_change", ch, int32(oldClass), int32(newClass))
			}
		}
	}
	ch.listener.OnShowEffect(ch, response.EffectTypeClassChange)
}

// openFourthClassSkills puts the class's visible fourth class skills on the skill list at level 0 and their WZ master
// level, from the third advancement on, as the original baseSkills did.
func (ch *Character) openFourthClassSkills(class constant.ClassType) {
	if class%10 < 1 || ch.gameWorld == nil {
		return
	}
	for id, model := range ch.gameWorld.GetResources().Skills {
		if model.Class != class || model.Invisible || !model.FourthClass || model.MasterLevel <= 0 {
			continue
		}
		if entry := ch.skills.Learned(id); entry != nil && (entry.Level() > 0 || entry.masterLevel > 0) {
			continue
		}
		ch.skills.Register(id, NewSkillEntry(ch, model, 0, model.MasterLevel))
	}
}

// ExpandInventory adds slots to a tab, up to 96.
func (ch *Character) ExpandInventory(invType constant.InventoryType, amount uint8) {
	tab := ch.inventory.tabs[invType]
	if tab == nil {
		return
	}
	limit := uint8(min(int(tab.slotLimit)+int(amount), inventorySlotMax))
	if limit == tab.slotLimit {
		return
	}
	tab.slotLimit = limit
	ch.listener.OnInventorySlotLimitChanged(ch, invType, limit)
}

// ResetInventorySlots puts every tab back to the slots a new character starts with.
func (ch *Character) ResetInventorySlots() {
	for invType, tab := range ch.inventory.tabs {
		if tab == nil || tab.slotLimit == inventorySlotDefault {
			continue
		}
		tab.slotLimit = inventorySlotDefault
		ch.listener.OnInventorySlotLimitChanged(ch, invType, inventorySlotDefault)
	}
}

func (ch *Character) grantClassChangePoints(newClass constant.ClassType) {
	if newClass == constant.ClassBeginner || newClass == constant.ClassNoblesse {
		return
	}

	ch.points.sp++

	if newClass%10 >= 2 {
		ch.points.sp += 2
	}

	if newClass%10 >= 1 && ch.level >= 70 {
		ch.points.ap += 5
	}

	if newClass%100 == 0 {
		minLevel := uint8(10)
		if newClass == constant.ClassMagician {
			minLevel = 8
		}

		if ch.level > minLevel {
			spToGrant := uint16(3 * (int(ch.level) - int(minLevel)))
			ch.points.sp += spToGrant
		}
	}
}
