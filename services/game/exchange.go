package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/util"
)

type ExchangeResult int

const (
	ExchangeOK ExchangeResult = iota
	ExchangeLackCost
	ExchangeLackCapacity
)

type ExchangeSkill struct {
	SkillID     uint32
	Level       int
	MasterLevel int
}

type exchangeObject struct {
	item  Item
	count uint16
}

type ExchangeSide struct {
	Items      map[uint32]uint16
	Period     map[uint32]time.Duration
	Meso       int32
	Exp        uint32
	Population int32
	SkillPoint uint16
	Skills     []ExchangeSkill
	Randomize  bool
	Bonus      map[string]int16
	objects    []exchangeObject
}

type ExchangeSpec struct {
	Cost   ExchangeSide
	Reward ExchangeSide
}

// include adds count of an existing item; its kind and expiration decide its slots, as when it is added.
func (side *ExchangeSide) include(item Item, count uint16) {
	side.objects = append(side.objects, exchangeObject{item: item, count: count})
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
	return side.Meso <= 0 && side.Exp == 0 && side.Population <= 0 && side.SkillPoint == 0 && exchangeItemsEmpty(side.Items) && len(side.objects) == 0 && len(side.Skills) == 0
}

func (spec ExchangeSpec) Valid(ch *Character, escrow *Escrow) ExchangeResult {
	if ch == nil {
		return ExchangeLackCost
	}
	if spec.Cost.isEmpty() && spec.Reward.isEmpty() {
		return ExchangeOK
	}
	if ch.Trading() && spec.Reward.Meso > 0 {
		return ExchangeLackCapacity
	}

	if result := ch.inventory.validateMeso(spec.Cost.Meso, spec.Reward.Meso); result != ExchangeOK {
		return result
	}
	if result := ch.validateExpExchange(spec.Cost.Exp); result != ExchangeOK {
		return result
	}
	if result := ch.stats.validatePopulation(spec.Cost.Population, spec.Reward.Population); result != ExchangeOK {
		return result
	}
	if result := spec.validateItems(ch, escrow); result != ExchangeOK {
		return result
	}

	// A trade passes no escrow: what the receiver offered leaves with the trade.
	holdsUnique := ch.inventory.hasUniqueItem
	if escrow == nil {
		holdsUnique = ch.inventory.carriesUniqueItem
	}
	rewards := make(map[uint32]uint16, len(spec.Reward.Items)+len(spec.Reward.objects))
	for id, count := range spec.Reward.Items {
		rewards[id] += count
	}
	for _, obj := range spec.Reward.objects {
		rewards[obj.item.GetModel().GetID()] += obj.count
	}
	for id, count := range rewards {
		if count == 0 {
			continue
		}
		model := ch.itemModel(id)
		if model == nil {
			return ExchangeLackCapacity
		}
		traits := model.GetTraits()
		if traits&constant.ItemTraitUnique == 0 || traits&constant.ItemTraitPartyQuest != 0 {
			continue
		}
		if count > 1 || (spec.Cost.Items[id] == 0 && holdsUnique(model)) {
			return ExchangeLackCapacity
		}
	}

	return ExchangeOK
}

// validateItems runs the cost and the reward through a slot plan of each tab, the same plan a real add uses.
func (spec ExchangeSpec) validateItems(ch *Character, escrow *Escrow) ExchangeResult {
	type tabPlan struct {
		plan *slotPlan
		ok   bool
	}
	plans := make(map[constant.InventoryType]tabPlan)
	planOf := func(invType constant.InventoryType, reward bool) (*slotPlan, ExchangeResult) {
		tp, made := plans[invType]
		if !made {
			if tab := ch.inventory.tabs[invType]; tab != nil {
				tp.plan, tp.ok = newSlotPlan(tab, invType, escrow)
			}
			plans[invType] = tp
		}
		switch {
		case tp.plan == nil && reward:
			return nil, ExchangeLackCapacity
		case tp.plan == nil:
			return nil, ExchangeLackCost
		case !tp.ok && reward:
			return nil, ExchangeLackCapacity
		}
		return tp.plan, ExchangeOK
	}

	for id, count := range spec.Cost.Items {
		if count == 0 {
			continue
		}
		model := ch.itemModel(id)
		if model == nil {
			return ExchangeLackCost
		}
		p, result := planOf(model.GetInventoryType(), false)
		if result != ExchangeOK {
			return result
		}
		if !p.take(id, count) {
			return ExchangeLackCost
		}
	}

	for id, count := range spec.Reward.Items {
		if count == 0 {
			continue
		}
		model := ch.itemModel(id)
		if model == nil {
			return ExchangeLackCapacity
		}
		p, result := planOf(model.GetInventoryType(), true)
		if result != ExchangeOK {
			return result
		}
		expiration := util.TimeMax
		if period := spec.Reward.Period[id]; period > 0 {
			expiration = clock.Now().Add(period)
		}
		if _, fits := p.place(model, count, expiration); !fits {
			return ExchangeLackCapacity
		}
	}
	for _, obj := range spec.Reward.objects {
		model := obj.item.GetModel()
		p, result := planOf(model.GetInventoryType(), true)
		if result != ExchangeOK {
			return result
		}
		if _, fits := p.place(model, obj.count, obj.item.GetExpiration()); !fits {
			return ExchangeLackCapacity
		}
	}
	return ExchangeOK
}

func (side ExchangeSide) ValidCost(ch *Character) ExchangeResult {
	if ch == nil {
		return ExchangeLackCost
	}
	if side.isEmpty() {
		return ExchangeOK
	}
	if side.Meso > 0 && ch.inventory.meso < side.Meso {
		return ExchangeLackCost
	}
	if side.Population > 0 && side.Population > int32(ch.stats.population) {
		return ExchangeLackCost
	}
	if side.Exp > 0 && ch.exp < side.Exp {
		return ExchangeLackCost
	}
	modelOf := ch.itemModel
	for id, count := range side.Items {
		if count == 0 {
			continue
		}
		if modelOf(id) == nil {
			return ExchangeLackCost
		}
		if !ch.inventory.HasItemCount(id, count) {
			return ExchangeLackCost
		}
	}
	return ExchangeOK
}
