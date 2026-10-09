package entity

func (ch *Character) NeedsQuestItem(questID uint32, itemID uint32) bool {
	qp := ch.Quests.Get(questID)
	if !qp.IsStarted() {
		return false
	}
	for _, act := range qp.Wz.Complete.Actions.Item {
		if act.ItemID == itemID && act.Count < 0 {
			return !ch.Inventory.HasItemCount(itemID, uint16(-act.Count))
		}
	}
	return true
}
