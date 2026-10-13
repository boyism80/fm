package entity

import "github.com/boyism80/fm/services/game/constant"

func (ch *Character) CloseChalkboard() {
	if ch.Chalkboard == "" || ch.GetMap() == nil {
		return
	}
	ch.Chalkboard = ""
	ch.Listener.OnChalkboard(ch)
}

func (ch *Character) SetChalkboard(text string) bool {
	if ch.GetMap() == nil || !ch.IsAlive() || ch.MiniRoom != nil || !ClientText(text).Valid(40) {
		return false
	}
	if ch.Chalkboard == text {
		return false
	}
	ch.Chalkboard = text
	ch.Listener.OnChalkboard(ch)
	return true
}

func (ch *Character) PlaceMessageBox(itemID uint32, text string) bool {
	m := ch.GetMap()
	if m == nil || !ch.IsAlive() || m.TemplateID()/10000000 == 91 || !ClientText(text).ValidLines(messageBoxLines, messageBoxLineLimit) {
		return false
	}
	box := &MessageBox{ItemID: itemID, Owner: ch.GetName(), Text: text}
	box.Position = ch.Position
	if err := m.AddMessageBox(box); err != nil {
		ch.Listener.OnMessageBoxFailed(ch)
		return false
	}
	return true
}

func (ch *Character) Express(expression uint32) {
	if ch.GetMap() == nil || !ch.IsAlive() || ch.Spectating() || expression == 0 || expression > 22 {
		return
	}
	if expression > 7 && ch.Inventory.Tabs[constant.InventoryTypeCash].Find(5159992+expression) == nil {
		return
	}
	ch.Listener.OnFaceExpression(ch, expression)
}

func (ch *Character) SetItemEffect(itemID uint32) {
	if ch.GetMap() == nil || !ch.IsAlive() || ch.Spectating() || ch.ItemEffectID == itemID {
		return
	}
	if itemID != 0 {
		item := ch.Inventory.Tabs[constant.InventoryTypeCash].Find(itemID)
		if item == nil || item.GetModel().GetTraits()&constant.ItemTraitItemEffect == 0 {
			return
		}
	}
	ch.ItemEffectID = itemID
	ch.Listener.OnItemEffect(ch)
}

func (ch *Character) clearItemEffect() {
	if ch.ItemEffectID == 0 {
		return
	}
	ch.ItemEffectID = 0
	if ch.GetMap() != nil {
		ch.Listener.OnItemEffect(ch)
	}
}

func (ch *Character) ShowChargeSkillEffect(skillID uint32, level, flags, speed uint8) {
	if ch.GetMap() == nil || !ch.IsAlive() || ch.Spectating() {
		return
	}
	skill := ch.Skills.Learned(skillID)
	if skill == nil || skill.Level() != int(level) || !skill.Wz.Charge {
		return
	}
	ch.Listener.OnChargeSkillEffect(ch, skillID, level, flags, speed)
}

func (ch *Character) CancelBuff(sourceID int32) {
	if sourceID > 0 {
		skill := ch.Skills.Learned(uint32(sourceID))
		if skill != nil && skill.Wz.Charge {
			ch.Listener.OnChargeSkillCancelled(ch, uint32(sourceID))
			return
		}
	}
	ch.Buffs.CancelBySource(sourceID)
}
