package entity

import (
	"math"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) itemModel(itemID uint32) wz.Item {
	if ch == nil || ch.GameWorld == nil {
		return nil
	}
	resources := ch.GameWorld.GetResources()
	if resources == nil {
		return nil
	}
	return resources.Items[itemID]
}

func (ch *Character) validateMesoFlow(costMeso, rewardMeso int32) FlowResult {
	if ch == nil {
		return FlowLackCost
	}
	if costMeso < 0 || rewardMeso < 0 {
		return FlowLackCost
	}
	if costMeso == 0 && rewardMeso == 0 {
		return FlowOK
	}
	if ch.Meso < costMeso {
		return FlowLackCost
	}
	after := ch.Meso - costMeso
	if rewardMeso > math.MaxInt32-after {
		return FlowLackCapacity
	}
	return FlowOK
}

func (ch *Character) validatePopulationFlow(costPopulation, rewardPopulation int32) FlowResult {
	if ch == nil {
		return FlowLackCost
	}
	if costPopulation < 0 || rewardPopulation < 0 {
		return FlowLackCost
	}
	if costPopulation == 0 && rewardPopulation == 0 {
		return FlowOK
	}
	current := uint32(ch.population)
	if uint32(costPopulation) > current {
		return FlowLackCost
	}
	after := current - uint32(costPopulation)
	if uint32(rewardPopulation) > 65535-after {
		return FlowLackCapacity
	}
	return FlowOK
}

func (ch *Character) validateExpFlow(costExp, rewardExp uint32) FlowResult {
	if ch == nil {
		return FlowLackCost
	}
	if costExp > 0 && ch.exp < costExp {
		return FlowLackCost
	}
	if rewardExp == 0 {
		return FlowOK
	}
	if rewardExp > ch.remainingExpToMaxLevel() {
		return FlowLackCapacity
	}
	return FlowOK
}

func (ch *Character) ValidateFlow(spec FlowSpec) FlowResult {
	if ch == nil {
		return FlowLackCost
	}
	if spec.Cost.isEmpty() && spec.Reward.isEmpty() {
		return FlowOK
	}

	if result := ch.validateMesoFlow(spec.Cost.Meso, spec.Reward.Meso); result != FlowOK {
		return result
	}
	if result := ch.validateExpFlow(spec.Cost.Exp, spec.Reward.Exp); result != FlowOK {
		return result
	}
	if result := ch.validatePopulationFlow(spec.Cost.Population, spec.Reward.Population); result != FlowOK {
		return result
	}

	modelOf := ch.itemModel
	costByInv := groupFlowItemsByInventory(spec.Cost.Items, modelOf)
	rewardByInv := groupFlowItemsByInventory(spec.Reward.Items, modelOf)

	invTypes := make(map[constant.InventoryType]struct{})
	for invType := range costByInv {
		invTypes[invType] = struct{}{}
	}
	for invType := range rewardByInv {
		invTypes[invType] = struct{}{}
	}

	for invType := range invTypes {
		var inv *Inventory
		if ch.Inventory != nil {
			inv = ch.Inventory[invType]
		}
		cost := costByInv[invType]
		reward := rewardByInv[invType]
		if result := inv.validateItemFlow(cost, reward, modelOf); result != FlowOK {
			return result
		}
	}

	if !flowItemsEmpty(spec.Cost.Items) {
		for id, count := range spec.Cost.Items {
			if count == 0 {
				continue
			}
			if modelOf(id) == nil {
				return FlowLackCost
			}
		}
	}
	if !flowItemsEmpty(spec.Reward.Items) {
		for id, count := range spec.Reward.Items {
			if count == 0 {
				continue
			}
			if modelOf(id) == nil {
				return FlowLackCapacity
			}
		}
	}

	return FlowOK
}

func (ch *Character) Exchange(spec FlowSpec) FlowResult {
	result := ch.ValidateFlow(spec)
	if result != FlowOK {
		return result
	}
	ch.applyFlowCost(spec.Cost)
	ch.applyFlowReward(spec.Reward)
	return FlowOK
}

func (ch *Character) applyFlowCost(side FlowSide) {
	if ch == nil || side.isEmpty() {
		return
	}
	if side.Meso > 0 {
		ch.removeMesoUnchecked(side.Meso)
	}
	if side.Population > 0 {
		ch.losePopulationUnchecked(side.Population)
	}
	for id, count := range side.Items {
		if count == 0 {
			continue
		}
		ch.removeByItemIDCountUnchecked(id, count)
	}
}

func (ch *Character) applyFlowReward(side FlowSide) {
	if ch == nil || side.isEmpty() {
		return
	}
	if side.Meso > 0 {
		ch.gainMesoUnchecked(side.Meso)
	}
	for id, count := range side.Items {
		if count == 0 {
			continue
		}
		if ch.GameWorld == nil {
			continue
		}
		item, err := NewItem(id, count, ch.GameWorld)
		if err != nil {
			continue
		}
		ch.applyAddItem(item, true)
	}
	if side.Exp > 0 {
		ch.addExpUnchecked(side.Exp)
	}
	if side.Population > 0 {
		ch.gainPopulationUnchecked(side.Population)
	}
}
