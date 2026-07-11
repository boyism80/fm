package entity

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type ExchangeResult int

const (
	ExchangeOK ExchangeResult = iota
	ExchangeLackCost
	ExchangeLackCapacity
)

type ExchangeSide struct {
	Items      map[uint32]uint16
	Meso       int32
	Exp        uint32
	Population int32
}

type ExchangeSpec struct {
	Cost   ExchangeSide
	Reward ExchangeSide
}

func exchangeItemsEmpty(items map[uint32]uint16) bool {
	if len(items) == 0 {
		return true
	}
	for _, count := range items {
		if count > 0 {
			return false
		}
	}
	return true
}

func (side ExchangeSide) isEmpty() bool {
	return side.Meso <= 0 && side.Exp == 0 && side.Population <= 0 && exchangeItemsEmpty(side.Items)
}

func (spec ExchangeSpec) Valid(ch *Character) ExchangeResult {
	if ch == nil {
		return ExchangeLackCost
	}
	if spec.Cost.isEmpty() && spec.Reward.isEmpty() {
		return ExchangeOK
	}

	if result := ch.validateMesoExchange(spec.Cost.Meso, spec.Reward.Meso); result != ExchangeOK {
		return result
	}
	if result := ch.validateExpExchange(spec.Cost.Exp, spec.Reward.Exp); result != ExchangeOK {
		return result
	}
	if result := ch.validatePopulationExchange(spec.Cost.Population, spec.Reward.Population); result != ExchangeOK {
		return result
	}

	modelOf := ch.itemModel
	costByInv := groupExchangeItemsByInventory(spec.Cost.Items, modelOf)
	rewardByInv := groupExchangeItemsByInventory(spec.Reward.Items, modelOf)

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
		if result := inv.validateItemExchange(cost, reward, modelOf); result != ExchangeOK {
			return result
		}
	}

	if !exchangeItemsEmpty(spec.Cost.Items) {
		for id, count := range spec.Cost.Items {
			if count == 0 {
				continue
			}
			if modelOf(id) == nil {
				return ExchangeLackCost
			}
		}
	}
	if !exchangeItemsEmpty(spec.Reward.Items) {
		for id, count := range spec.Reward.Items {
			if count == 0 {
				continue
			}
			if modelOf(id) == nil {
				return ExchangeLackCapacity
			}
		}
	}

	return ExchangeOK
}

func groupExchangeItemsByInventory(items map[uint32]uint16, modelOf func(uint32) wz.Item) map[constant.InventoryType]map[uint32]uint16 {
	if len(items) == 0 {
		return nil
	}
	out := make(map[constant.InventoryType]map[uint32]uint16)
	for id, count := range items {
		if count == 0 {
			continue
		}
		if modelOf != nil && modelOf(id) == nil {
			continue
		}
		invType := constant.GetInventoryTypeByItemID(id)
		if out[invType] == nil {
			out[invType] = make(map[uint32]uint16)
		}
		out[invType][id] += count
	}
	return out
}
