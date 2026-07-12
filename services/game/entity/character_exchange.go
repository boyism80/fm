package entity

import (
	"math"
	"time"

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

func (ch *Character) validateMesoExchange(costMeso, rewardMeso int32) ExchangeResult {
	if ch == nil {
		return ExchangeLackCost
	}
	if costMeso < 0 || rewardMeso < 0 {
		return ExchangeLackCost
	}
	if costMeso == 0 && rewardMeso == 0 {
		return ExchangeOK
	}
	if ch.Meso < costMeso {
		return ExchangeLackCost
	}
	after := ch.Meso - costMeso
	if rewardMeso > math.MaxInt32-after {
		return ExchangeLackCapacity
	}
	return ExchangeOK
}

func (ch *Character) validatePopulationExchange(costPopulation, rewardPopulation int32) ExchangeResult {
	if ch == nil {
		return ExchangeLackCost
	}
	if costPopulation < 0 || rewardPopulation < 0 {
		return ExchangeLackCost
	}
	if costPopulation == 0 && rewardPopulation == 0 {
		return ExchangeOK
	}
	current := uint32(ch.population)
	if uint32(costPopulation) > current {
		return ExchangeLackCost
	}
	after := current - uint32(costPopulation)
	if uint32(rewardPopulation) > 65535-after {
		return ExchangeLackCapacity
	}
	return ExchangeOK
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
	result := spec.Valid(ch)
	if result != ExchangeOK {
		return result
	}
	ch.applyExchangeCost(spec.Cost)
	ch.applyExchangeReward(spec.Reward)
	return ExchangeOK
}

func (ch *Character) applyExchangeCost(side ExchangeSide) {
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
	for _, skill := range side.Skills {
		if skill.SkillID == 0 || ch.Skills == nil {
			continue
		}
		ch.Skills.Remove(skill.SkillID)
	}
}

func (ch *Character) applyExchangeReward(side ExchangeSide) {
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
	for _, skill := range side.Skills {
		ch.applyExchangeSkillGrant(skill)
	}
}

func (ch *Character) applyExchangeSkillGrant(skill ExchangeSkill) {
	if ch == nil || skill.SkillID == 0 || ch.Skills == nil {
		return
	}
	if skill.SkillID/10000 == 0 && !ch.IsBeginner() {
		return
	}
	if ch.GameWorld == nil {
		return
	}
	resources := ch.GameWorld.GetResources()
	if resources == nil {
		return
	}
	wzSkill := resources.GetSkill(skill.SkillID)
	if wzSkill == nil {
		return
	}
	targetLevel := exchangeSkillTargetLevel(0, skill.Level)
	targetMaster := exchangeSkillTargetMaster(0, skill.MasterLevel, wzSkill)
	entry := ch.Skills.Get(skill.SkillID)
	if entry == nil {
		entry = NewSkillEntry(ch, wzSkill, targetLevel, targetMaster)
		entry.Expiration = time.Time{}
		ch.Skills.Register(skill.SkillID, entry)
	} else {
		targetLevel = exchangeSkillTargetLevel(entry.Level(), skill.Level)
		targetMaster = exchangeSkillTargetMaster(entry.MasterLevel, skill.MasterLevel, wzSkill)
		if targetLevel != entry.Level() || targetMaster != entry.MasterLevel {
			entry.SetLevelAndMaster(targetLevel, targetMaster)
		}
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

func exchangeSkillTargetMaster(current int, rewardMaster int, wzSkill *wz.Skill) int {
	resolved := rewardMaster
	if resolved <= 0 && wzSkill != nil {
		if wzSkill.MasterLevel > 0 {
			resolved = wzSkill.MasterLevel
		} else if wzSkill.MaxLevel > 0 {
			resolved = wzSkill.MaxLevel
		}
	}
	if current > resolved {
		return current
	}
	return resolved
}
