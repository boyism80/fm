package entity

import (
	"math/rand"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) UseItemMaker(typ uint32, itemID uint32, catalyst bool, enhancers []uint32, slot int16) {
	if ch.GetHp() == 0 || ch.GameWorld == nil {
		return
	}
	switch typ {
	case 1, 2:
		ch.makeItem(itemID, catalyst, enhancers)
	case 3:
		ch.makeCrystal(itemID)
	case 4:
		ch.disassembleItem(slot, itemID)
	}
}

func (ch *Character) makerSkillLevel() int {
	skillID := uint32(constant.SkillMaker)
	if ch.Class >= 1000 && ch.Class < 2000 {
		skillID = uint32(constant.SkillMakerCygnus)
	}
	entry := ch.Skills.Get(skillID)
	if entry == nil {
		return 0
	}
	return entry.Level()
}

func (ch *Character) makeItem(itemID uint32, catalyst bool, enhancers []uint32) {
	resources := ch.GameWorld.GetResources()
	recipe := resources.ItemMake[itemID]
	if recipe == nil || int(ch.GetLevel()) < recipe.RequiredLevel || ch.makerSkillLevel() < recipe.RequiredSkillLevel {
		return
	}
	if recipe.Class != 0 {
		category := int(ch.Class) / 100 % 10
		if category < 1 || recipe.Class&(1<<(category-1)) == 0 {
			return
		}
	}
	if recipe.RequiredItem != 0 && !ch.Inventory.HasItem(recipe.RequiredItem) {
		return
	}
	if recipe.RequiredEquip != 0 {
		owned := ch.Inventory.HasItem(recipe.RequiredEquip)
		for _, equipment := range ch.Inventory.Equipped {
			if equipment.GetModel().GetID() == recipe.RequiredEquip {
				owned = true
			}
		}
		if !owned {
			return
		}
	}
	if len(enhancers) > recipe.EnhancerCount {
		return
	}

	cost := ExchangeSide{Items: make(map[uint32]uint16)}
	for _, ingredient := range recipe.Ingredients {
		if ingredient.Count > ^uint16(0)-cost.Items[ingredient.ItemID] {
			return
		}
		cost.Items[ingredient.ItemID] += ingredient.Count
	}
	rate := 1
	if catalyst {
		if recipe.Catalyst == 0 {
			return
		}
		cost.Items[recipe.Catalyst]++
		rate += 5
	}
	bonus := make(map[string]int16)
	for _, enhancerID := range enhancers {
		model := resources.Items[enhancerID]
		if model == nil || len(model.GetMakerBonus()) == 0 {
			return
		}
		cost.Items[enhancerID]++
		rate += 3 * (int(enhancerID%10) + 1)
		for stat, value := range model.GetMakerBonus() {
			bonus[stat] += value
		}
	}

	// The client shows and checks this price before sending the request.
	meso := int(recipe.Meso) + max(int(recipe.Meso)/10*rate, 0)
	cost.Meso = int32(meso - meso%1000)

	rewardID := recipe.ItemID
	rewardCount := recipe.Count
	if len(recipe.RandomRewards) > 0 {
		total := 0
		for _, reward := range recipe.RandomRewards {
			total += reward.Weight
		}
		roll := rand.Intn(total)
		for _, reward := range recipe.RandomRewards {
			roll -= reward.Weight
			if roll < 0 {
				rewardID = reward.ItemID
				rewardCount = reward.Count
				break
			}
		}
	}

	destroyed := catalyst && rand.Intn(10) == 0
	spec := ExchangeSpec{Cost: cost}
	if !destroyed {
		spec.Reward.Items = map[uint32]uint16{rewardID: rewardCount}
	}
	if spec.Valid(ch) != ExchangeOK {
		return
	}

	var made Item
	if !destroyed {
		item, err := NewItem(rewardID, rewardCount, ch.GameWorld)
		if err != nil {
			return
		}
		if equipment, ok := item.(Equipment); ok {
			if catalyst {
				model, ok := equipment.GetModel().(wz.Equipment)
				if !ok {
					return
				}
				equipment.GetEquipmentCore().RandomizeStats(model)
			}
			equipment.GetEquipmentCore().AddBonus(bonus)
		}
		made = item
	}

	ch.payExchangeCost(spec.Cost)
	if made == nil {
		ch.Listener.OnItemMakerDestroyed(ch)
		return
	}
	ch.Inventory.addItemUnchecked(made, true)
	ch.Listener.OnShowSelfItemMakerSuccessEffect(ch)
}

func (ch *Character) makeCrystal(itemID uint32) {
	model := ch.itemModel(itemID)
	if model == nil {
		return
	}
	crystalID := ch.makerCrystalID(model.GetMakerLevel())
	if crystalID == 0 {
		return
	}
	spec := ExchangeSpec{
		Cost:   ExchangeSide{Items: map[uint32]uint16{itemID: 100}},
		Reward: ExchangeSide{Items: map[uint32]uint16{crystalID: 1}},
	}
	if ch.Exchange(spec) != ExchangeOK {
		return
	}
	ch.Listener.OnShowSelfItemMakerSuccessEffect(ch)
}

func (ch *Character) disassembleItem(slot int16, itemID uint32) {
	item := ch.Inventory.GetItem(constant.InventoryTypeEquipment, slot)
	equipment, ok := item.(Equipment)
	if !ok || item.GetModel().GetID() != itemID || item.GetModel().GetTraits()&(constant.ItemTraitTradeBlock|constant.ItemTraitAccountSharable) != 0 {
		return
	}
	model, ok := item.GetModel().(wz.Equipment)
	if !ok {
		return
	}
	crystalID := ch.makerCrystalID(int(model.GetRequired().Level))
	if crystalID == 0 {
		return
	}

	quality := equipment.GetEquipmentCore().BonusStats.quality()
	count := 3 + rand.Intn(5)
	if model.GetEquipmentType() == constant.EquipmentTypeWeapon || model.GetEquipmentType() == constant.EquipmentTypeLongcoat {
		count = 5 + rand.Intn(7)
	}
	if quality > 0 {
		count += rand.Intn(quality) + 1
	} else if quality < 0 {
		count -= rand.Intn(-quality + 1)
	}
	count = max(count, 1)

	// The client shows and checks this price before sending the request.
	base := 150
	switch quality {
	case -1:
		base = 100
	case 1:
		base = 200
	case 2:
		base = 250
	case 3:
		base = 300
	}
	meso := base * (item.GetModel().GetPrice() / 10) / 100
	spec := ExchangeSpec{
		Cost:   ExchangeSide{Meso: int32(meso - meso%1000)},
		Reward: ExchangeSide{Items: map[uint32]uint16{crystalID: uint16(count)}},
	}
	if spec.Valid(ch) != ExchangeOK {
		return
	}
	if !ch.Inventory.RemoveItem(constant.InventoryTypeEquipment, slot, 1) {
		return
	}
	ch.payExchangeCost(spec.Cost)
	ch.grantExchangeReward(spec.Reward)
	ch.Listener.OnShowSelfItemMakerSuccessEffect(ch)
}

func (ch *Character) makerCrystalID(level int) uint32 {
	if level < 31 {
		return 0
	}
	if level <= 50 {
		return 4260000
	}
	return 4260000 + uint32(min((level-41)/10, 8))
}
