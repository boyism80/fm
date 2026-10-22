package game

import (
	"math/rand"
	"slices"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

func (inv *Inventory) Make(typ uint32, itemID uint32, catalyst bool, enhancers []uint32, slot int16) {
	if inv.owner.GetHp() == 0 || inv.owner.gameWorld == nil {
		return
	}
	switch typ {
	case 1, 2:
		inv.makeItem(itemID, catalyst, enhancers)
	case 3:
		inv.makeCrystal(itemID)
	case 4:
		inv.disassemble(slot, itemID)
	}
}

func (inv *Inventory) makeItem(itemID uint32, catalyst bool, enhancers []uint32) {
	resources := inv.owner.gameWorld.GetResources()
	recipe := resources.ItemMake[itemID]
	if recipe == nil || int(inv.owner.GetLevel()) < recipe.RequiredLevel || inv.owner.skills.makerLevel() < recipe.RequiredSkillLevel {
		return
	}
	if recipe.Classes != nil && !slices.ContainsFunc(recipe.Classes, inv.owner.class.Of) {
		return
	}
	if recipe.RequiredItem != 0 && !inv.HasItem(recipe.RequiredItem) {
		return
	}
	if recipe.RequiredEquip != 0 {
		owned := inv.HasItem(recipe.RequiredEquip)
		for _, equipment := range inv.equipped {
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
	gems := make([]wz.Item, 0, len(enhancers))
	for _, enhancerID := range enhancers {
		model := resources.Items[enhancerID]
		if model == nil || model.GetTraits()&constant.ItemTraitMakerGem == 0 {
			return
		}
		cost.Items[enhancerID]++
		rate += 3 * (int(enhancerID%10) + 1)
		gems = append(gems, model)
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
	if spec.Valid(inv.owner, inv.owner.escrow) != ExchangeOK {
		return
	}

	var made Item
	if !destroyed {
		item, err := NewItem(rewardID, rewardCount, inv.owner.gameWorld)
		if err != nil {
			return
		}
		if equipment, ok := item.(Equipment); ok {
			model, ok := equipment.GetModel().(wz.Equipment)
			if !ok {
				return
			}
			if catalyst {
				equipment.GetEquipmentCore().RandomizeStats(model)
			}
			for _, gem := range gems {
				bonus := gem.GetMakerBonus()
				equipment.GetEquipmentCore().AddBonus(bonus)
				equipment.GetEquipmentCore().RandomizeMakerOption(model, bonus["randOption"], bonus["randStat"])
			}
		}
		made = item
	}

	inv.owner.payExchangeCost(spec.Cost)
	if made == nil {
		inv.owner.listener.OnItemMakerDestroyed(inv.owner)
		return
	}
	inv.addItemUnchecked(made, true)
	inv.owner.listener.OnShowSelfItemMakerSuccessEffect(inv.owner)
}

func (inv *Inventory) makeCrystal(itemID uint32) {
	model := inv.owner.itemModel(itemID)
	if model == nil {
		return
	}
	crystalID := inv.makerCrystalID(model.GetMakerLevel())
	if crystalID == 0 {
		return
	}
	spec := ExchangeSpec{
		Cost:   ExchangeSide{Items: map[uint32]uint16{itemID: 100}},
		Reward: ExchangeSide{Items: map[uint32]uint16{crystalID: 1}},
	}
	if inv.owner.Exchange(spec) != ExchangeOK {
		return
	}
	inv.owner.listener.OnShowSelfItemMakerSuccessEffect(inv.owner)
}

func (inv *Inventory) disassemble(slot int16, itemID uint32) {
	item := inv.GetItem(constant.InventoryTypeEquipment, slot)
	equipment, ok := item.(Equipment)
	if !ok || item.GetModel().GetID() != itemID || item.GetModel().GetTraits()&(constant.ItemTraitTradeBlock|constant.ItemTraitAccountSharable) != 0 {
		return
	}
	model, ok := item.GetModel().(wz.Equipment)
	if !ok {
		return
	}
	crystalID := inv.makerCrystalID(int(model.GetRequired().Level))
	if crystalID == 0 {
		return
	}

	quality := equipment.GetEquipmentCore().bonusStats.quality()
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
	if spec.Valid(inv.owner, inv.owner.escrow) != ExchangeOK {
		return
	}
	if !inv.RemoveItem(constant.InventoryTypeEquipment, slot, 1) {
		return
	}
	inv.owner.payExchangeCost(spec.Cost)
	inv.owner.grantExchangeReward(spec.Reward)
	inv.owner.listener.OnShowSelfItemMakerSuccessEffect(inv.owner)
}

func (inv *Inventory) makerCrystalID(level int) uint32 {
	if level < 31 {
		return 0
	}
	if level <= 50 {
		return 4260000
	}
	return 4260000 + uint32(min((level-41)/10, 8))
}
