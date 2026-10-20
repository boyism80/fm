package game

import "github.com/boyism80/fm/maple/constant"

const chatLimit = 70

func (ch *Character) Chat(message string, dontRecordHistory bool) {
	if ch.GetMap() == nil || ch.Forbidden("chat") {
		return
	}
	if !ch.HasRoleAtLeast(constant.RoleAdmin) && !ClientText(message).Valid(chatLimit) {
		return
	}
	ch.listener.OnChat(ch, message, false, dontRecordHistory)
}

// Forbidden is an event map rule; game masters running the event are not bound by it.
func (ch *Character) Forbidden(action string) bool {
	return !ch.HasRoleAtLeast(constant.RoleAdmin) && ch.GetMap().Forbids(action)
}

func (ch *Character) CloseChalkboard() {
	if ch.chalkboard == "" || ch.GetMap() == nil {
		return
	}
	ch.chalkboard = ""
	ch.listener.OnChalkboard(ch)
}

func (ch *Character) SetChalkboard(text string) bool {
	if ch.GetMap() == nil || !ch.IsAlive() || ch.miniRoom != nil || !ClientText(text).Valid(40) {
		return false
	}
	if ch.chalkboard == text {
		return false
	}
	ch.chalkboard = text
	ch.listener.OnChalkboard(ch)
	return true
}

func (ch *Character) Express(expression uint32) {
	if ch.GetMap() == nil || !ch.IsAlive() || ch.Spectating() || expression == 0 || expression > 22 {
		return
	}
	if expression > 7 && ch.inventory.tabs[constant.InventoryTypeCash].Find(5159992+expression) == nil {
		return
	}
	ch.listener.OnFaceExpression(ch, expression)
}

func (ch *Character) SetItemEffect(itemID uint32) {
	if ch.GetMap() == nil || !ch.IsAlive() || ch.Spectating() || ch.itemEffectID == itemID {
		return
	}
	if itemID != 0 {
		item := ch.inventory.tabs[constant.InventoryTypeCash].Find(itemID)
		if item == nil || item.GetModel().GetTraits()&constant.ItemTraitItemEffect == 0 {
			return
		}
	}
	ch.itemEffectID = itemID
	ch.listener.OnItemEffect(ch)
}

func (ch *Character) clearItemEffect() {
	if ch.itemEffectID == 0 {
		return
	}
	ch.itemEffectID = 0
	if ch.GetMap() != nil {
		ch.listener.OnItemEffect(ch)
	}
}

func (ch *Character) ShowChargeSkillEffect(skillID uint32, level, flags, speed uint8) {
	if ch.GetMap() == nil || !ch.IsAlive() || ch.Spectating() {
		return
	}
	skill := ch.skills.Learned(skillID)
	if skill == nil || skill.Level() != int(level) || !skill.wz.Charge {
		return
	}
	ch.listener.OnChargeSkillEffect(ch, skillID, level, flags, speed)
}
