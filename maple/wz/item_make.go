package wz

import (
	"encoding/xml"
	"os"
	"strconv"

	"github.com/boyism80/fm/maple/constant"
)

var itemMakeClassBits = []struct {
	mask    int
	classes []constant.ClassType
}{
	{0x1, []constant.ClassType{constant.ClassWarrior, constant.ClassDawnWarrior1, constant.ClassAran2}},
	{0x2, []constant.ClassType{constant.ClassMagician, constant.ClassBlazeWizard1}},
	{0x4, []constant.ClassType{constant.ClassBowman, constant.ClassWindArcher1}},
	{0x8, []constant.ClassType{constant.ClassThief, constant.ClassNightWalker1}},
	{0x10, []constant.ClassType{constant.ClassPirate, constant.ClassThunderBreaker1}},
}

var makerBonusStats = map[string]string{
	"incPAD":   "watk",
	"incMAD":   "matk",
	"incACC":   "acc",
	"incEVA":   "avoid",
	"incSpeed": "speed",
	"incJump":  "jump",
	"incMaxHP": "hp",
	"incMaxMP": "mp",
	"incSTR":   "str",
	"incINT":   "int",
	"incLUK":   "luk",
	"incDEX":   "dex",
}

type ItemMakeIngredient struct {
	ItemID uint32
	Count  uint16
}

type ItemMakeReward struct {
	ItemID uint32
	Count  uint16
	Weight int
}

type ItemMakeRecipe struct {
	ItemID             uint32
	Classes            []constant.ClassType
	RequiredLevel      int
	RequiredSkillLevel int
	Count              uint16
	EnhancerCount      int
	Meso               int32
	Catalyst           uint32
	RequiredItem       uint32
	RequiredEquip      uint32
	Ingredients        []ItemMakeIngredient
	RandomRewards      []ItemMakeReward
}

func loadItemMake(path string) (map[uint32]*ItemMakeRecipe, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}
	recipes := make(map[uint32]*ItemMakeRecipe)
	for _, group := range root.Children {
		mask, err := strconv.Atoi(group.Name)
		if err != nil {
			continue
		}
		var classes []constant.ClassType
		for _, bit := range itemMakeClassBits {
			if mask&bit.mask != 0 {
				classes = append(classes, bit.classes...)
			}
		}
		for _, entry := range group.Children {
			id, err := strconv.Atoi(entry.Name)
			if err != nil || id <= 0 {
				continue
			}
			recipe := &ItemMakeRecipe{
				ItemID:             uint32(id),
				Classes:            classes,
				RequiredLevel:      entry.Int("reqLevel", 0),
				RequiredSkillLevel: entry.Int("reqSkillLevel", 0),
				Count:              uint16(max(entry.Int("itemNum", 1), 1)),
				EnhancerCount:      entry.Int("tuc", 0),
				Meso:               int32(max(entry.Int("meso", 0), 0)),
				Catalyst:           uint32(max(entry.Int("catalyst", 0), 0)),
				RequiredItem:       uint32(max(entry.Int("reqItem", 0), 0)),
				RequiredEquip:      uint32(max(entry.Int("reqEquip", 0), 0)),
			}
			if ingredients := entry.find("recipe"); ingredients != nil {
				for _, ingredient := range ingredients.Children {
					itemID := ingredient.Int("item", 0)
					count := ingredient.Int("count", 0)
					if itemID > 0 && count > 0 && count <= 65535 {
						recipe.Ingredients = append(recipe.Ingredients, ItemMakeIngredient{ItemID: uint32(itemID), Count: uint16(count)})
					}
				}
			}
			if rewards := entry.find("randomReward"); rewards != nil {
				for _, reward := range rewards.Children {
					itemID := reward.Int("item", 0)
					count := reward.Int("itemNum", 0)
					weight := reward.Int("prob", 0)
					if itemID > 0 && count > 0 && count <= 65535 && weight > 0 {
						recipe.RandomRewards = append(recipe.RandomRewards, ItemMakeReward{ItemID: uint32(itemID), Count: uint16(count), Weight: weight})
					}
				}
			}
			recipes[recipe.ItemID] = recipe
		}
	}
	return recipes, nil
}
