package game

import (
	"slices"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

func (qc *Quests) OpenPot(questID uint32) {
	qp := qc.Get(questID)
	if !qp.IsStarted() || qp.statusRecord.IsEmpty() {
		return
	}
	qc.owner.listener.OnQuestProgress(qc.owner, qp)
}

func (qc *Quests) FeedPot(slot int16, itemID uint32, questID uint32) {
	item := qc.owner.inventory.GetItem(constant.InventoryTypeETC, slot)
	if item == nil || item.GetModel().GetID() != itemID {
		return
	}
	feed, ok := item.GetModel().(*wz.MiscItem)
	if !ok || feed.Exp == 0 {
		return
	}

	var pot *wz.MiscItem
	etc := qc.owner.inventory.tabs[constant.InventoryTypeETC]
	for s := int16(1); s <= int16(etc.slotLimit); s++ {
		candidate := etc.Get(uint8(s))
		if candidate == nil {
			continue
		}
		model, ok := candidate.GetModel().(*wz.MiscItem)
		if ok && model.QuestID == questID && model.MaxExp > 0 && slices.Contains(model.ConsumeItems, itemID) {
			pot = model
			break
		}
	}
	if pot == nil {
		return
	}

	qp := qc.Get(questID)
	if !qp.IsStarted() {
		return
	}
	exp, err := qp.statusRecord.AsInt()
	if err != nil || exp >= pot.MaxExp {
		return
	}
	if !qc.owner.inventory.RemoveItem(constant.InventoryTypeETC, slot, 1) {
		return
	}

	exp = min(exp+feed.Exp, pot.MaxExp)
	qp.statusRecord.WriteInt(exp)
	qc.owner.listener.OnQuestProgress(qc.owner, qp)
	if exp < pot.MaxExp {
		return
	}

	qc.ForEach(func(id uint32, dependent *Quest) {
		if uint32(dependent.wz.Complete.Requirements.InfoNumber) != questID {
			return
		}
		if dependent.CanComplete(qc.owner, QuestPhaseOpts{}) == nil {
			qc.owner.listener.OnShowQuestCompletion(qc.owner, id)
		}
	})
}
