package game

import (
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/wz"
)

func (ch *Character) itemModel(itemID uint32) wz.Item {
	if ch == nil || ch.gameWorld == nil {
		return nil
	}
	resources := ch.gameWorld.GetResources()
	if resources == nil {
		return nil
	}
	return resources.Items[itemID]
}

func (ch *Character) validateExpExchange(costExp, rewardExp uint32) ExchangeResult {
	if ch == nil {
		return ExchangeLackCost
	}
	if costExp > 0 && ch.exp < costExp {
		return ExchangeLackCost
	}
	if rewardExp == 0 {
		return ExchangeOK
	}
	if rewardExp > ch.remainingExpToMaxLevel() {
		return ExchangeLackCapacity
	}
	return ExchangeOK
}

func (ch *Character) Exchange(spec ExchangeSpec) ExchangeResult {
	result := spec.Valid(ch, ch.escrow)
	if result != ExchangeOK {
		return result
	}
	ch.payExchangeCost(spec.Cost)
	ch.grantExchangeReward(spec.Reward)
	return ExchangeOK
}

func (ch *Character) payExchangeCost(side ExchangeSide) {
	if side.isEmpty() {
		return
	}
	if side.Meso > 0 {
		ch.inventory.removeMesoUnchecked(side.Meso)
	}
	if side.Population > 0 {
		ch.stats.AddPopulation(-side.Population)
	}
	for id, count := range side.Items {
		if count == 0 {
			continue
		}
		ch.inventory.removeByItemIDCountUnchecked(id, count)
	}
	for _, skill := range side.Skills {
		if skill.SkillID == 0 {
			continue
		}
		ch.skills.Remove(skill.SkillID)
	}
}

func (ch *Character) grantExchangeReward(side ExchangeSide) {
	if side.isEmpty() {
		return
	}
	if side.Meso > 0 {
		ch.inventory.gainMesoUnchecked(side.Meso)
	}
	for id, count := range side.Items {
		if count == 0 {
			continue
		}
		if ch.gameWorld == nil {
			continue
		}
		item, err := NewItem(id, count, ch.gameWorld)
		if err != nil {
			continue
		}
		if eq, ok := item.(Equipment); ok {
			if side.Randomize {
				if em, ok := eq.GetModel().(wz.Equipment); ok {
					eq.GetEquipmentCore().RandomizeStats(em)
				}
			}
			eq.GetEquipmentCore().SetBonus(side.Bonus)
		}
		if period := side.Period[id]; period > 0 {
			item.SetExpiration(clock.Now().Add(period))
		}
		ch.inventory.addItemUnchecked(item, true)
	}
	if side.Exp > 0 {
		ch.addExpUnchecked(side.Exp)
	}
	if side.Population > 0 {
		ch.stats.AddPopulation(side.Population)
	}
	if side.SkillPoint > 0 {
		ch.points.SetSP(ch.points.sp+side.SkillPoint, true)
	}
	for _, skill := range side.Skills {
		ch.skills.grant(skill)
	}
}

func exchangeSkillTargetLevel(current int, rewardLevel int) int {
	if rewardLevel <= 0 {
		if current > 0 {
			return current
		}
		return 1
	}
	if current > rewardLevel {
		return current
	}
	return rewardLevel
}
