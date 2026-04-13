package wz

import (
	"encoding/xml"
	"os"
	"strconv"
)

func (sd *StringData) CountStrings() int {
	count := 0

	for _, regionData := range sd.MapStrings {
		count += len(regionData)
	}

	count += len(sd.MobStrings)
	count += len(sd.NpcStrings)
	count += len(sd.SkillStrings)
	count += len(sd.ItemCashStrings)
	count += len(sd.ItemConsumeStrings)
	count += len(sd.ItemEtcStrings)
	count += len(sd.ItemInsStrings)
	count += len(sd.ItemPetStrings)

	for _, categoryData := range sd.ItemEqpStrings {
		count += len(categoryData)
	}

	return count
}

func loadMapStrings(path string) (*map[string]map[uint32]map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	result := make(map[string]map[uint32]map[string]string)

	for _, regionNode := range root.Children {
		region := regionNode.Name
		if result[region] == nil {
			result[region] = make(map[uint32]map[string]string)
		}

		for _, mapNode := range regionNode.Children {
			mapId, err := strconv.ParseUint(mapNode.Name, 10, 32)
			if err != nil {
				continue
			}

			stringData := make(map[string]string)

			for _, field := range mapNode.Strings {
				stringData[field.Name] = field.Value
			}
			result[region][uint32(mapId)] = stringData
		}
	}

	return &result, nil
}

func loadMobStrings(path string) (*map[uint32]map[string]string, error) {
	return loadSimpleStrings(path)
}

func loadNpcStrings(path string) (*map[uint32]map[string]string, error) {
	return loadSimpleStrings(path)
}

func loadSkillStrings(path string) (*map[uint32]map[string]string, error) {
	return loadSimpleStrings(path)
}

func loadItemStrings(path string) (*map[uint32]map[string]string, error) {
	return loadSimpleStrings(path)
}

func loadEqpStrings(path string) (*map[string]map[uint32]map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	result := make(map[string]map[uint32]map[string]string)

	var processNode func(node *node, category string)
	processNode = func(n *node, category string) {

		if itemId, err := strconv.ParseUint(n.Name, 10, 32); err == nil {

			if result[category] == nil {
				result[category] = make(map[uint32]map[string]string)
			}

			stringData := make(map[string]string)

			for _, field := range n.Strings {
				stringData[field.Name] = field.Value
			}
			result[category][uint32(itemId)] = stringData
		} else {

			for _, child := range n.Children {

				childCategory := category
				if category == "" {
					childCategory = n.Name
				}
				processNode(&child, childCategory)
			}
		}
	}

	for _, rootChild := range root.Children {

		if rootChild.Name == "Eqp" {

			for _, categoryNode := range rootChild.Children {
				processNode(&categoryNode, categoryNode.Name)
			}
		} else {

			processNode(&rootChild, rootChild.Name)
		}
	}

	return &result, nil
}

func loadSimpleStrings(path string) (*map[uint32]map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	result := make(map[uint32]map[string]string)

	var processNode func(n *node)
	processNode = func(n *node) {

		id, err := strconv.ParseUint(n.Name, 10, 32)
		if err != nil {

			for i := range n.Children {
				processNode(&n.Children[i])
			}
			return
		}

		stringData := make(map[string]string)

		for _, field := range n.Strings {
			stringData[field.Name] = field.Value
		}
		result[uint32(id)] = stringData
	}

	for i := range root.Children {
		processNode(&root.Children[i])
	}

	return &result, nil
}

func getMapRegion(mapId uint32) string {
	if mapId < 100000000 {
		return "maple"
	} else if mapId < 200000000 {
		return "victoria"
	} else if mapId < 300000000 {
		return "ossyria"
	} else if mapId < 400000000 {
		return "3rd"
	} else if mapId >= 500000000 && mapId < 510000000 {
		return "thai"
	} else if mapId >= 555000000 && mapId < 556000000 {
		return "SG"
	} else if mapId >= 540000000 && mapId < 600000000 {
		return "singapore"
	}

	return "maple"
}

func (r *Resources) GetMapName(mapId uint32) (mapName, streetName string) {
	region := getMapRegion(mapId)
	if regionData, ok := r.Strings.MapStrings[region]; ok {
		if mapData, ok := regionData[mapId]; ok {
			mapName = mapData["mapName"]
			streetName = mapData["streetName"]
		}
	}
	return
}

func (r *Resources) GetMobName(mobId uint32) string {
	if mobData, ok := r.Strings.MobStrings[mobId]; ok {
		return mobData["name"]
	}
	return ""
}

func (r *Resources) GetNpcName(npcId uint32) string {
	if npcData, ok := r.Strings.NpcStrings[npcId]; ok {
		return npcData["name"]
	}
	return ""
}

func (r *Resources) GetSkillName(skillId uint32) string {
	if skillData, ok := r.Strings.SkillStrings[skillId]; ok {
		return skillData["name"]
	}
	return ""
}

func (r *Resources) GetItemName(itemId uint32) string {

	if itemId >= 5010000 {

		if itemData, ok := r.Strings.ItemCashStrings[itemId]; ok {
			return itemData["name"]
		}
	} else if itemId >= 2000000 && itemId < 3000000 {

		if itemData, ok := r.Strings.ItemConsumeStrings[itemId]; ok {
			return itemData["name"]
		}
	} else if itemId >= 4000000 && itemId < 5000000 {

		if itemData, ok := r.Strings.ItemEtcStrings[itemId]; ok {
			return itemData["name"]
		}
	} else if itemId >= 3000000 && itemId < 4000000 {

		if itemData, ok := r.Strings.ItemInsStrings[itemId]; ok {
			return itemData["name"]
		}
	} else if itemId >= 5000000 && itemId < 5010000 {

		if itemData, ok := r.Strings.ItemPetStrings[itemId]; ok {
			return itemData["name"]
		}
	} else {

		category := getItemCategory(itemId)
		if categoryData, ok := r.Strings.ItemEqpStrings[category]; ok {
			if itemData, ok := categoryData[itemId]; ok {
				return itemData["name"]
			}
		}
	}
	return ""
}

func getItemCategory(itemId uint32) string {

	if itemId >= 1132000 && itemId < 1183000 || (itemId >= 1010000 && itemId < 1040000) || (itemId >= 1122000 && itemId < 1123000) {
		return "Accessory"
	} else if itemId >= 1172000 && itemId < 1173000 {
		return "MonsterBook"
	} else if itemId >= 1662000 && itemId < 1680000 {
		return "Android"
	} else if itemId >= 1000000 && itemId < 1010000 {
		return "Cap"
	} else if itemId >= 1102000 && itemId < 1104000 {
		return "Cape"
	} else if itemId >= 1040000 && itemId < 1050000 {
		return "Coat"
	} else if itemId >= 20000 && itemId < 22000 {
		return "Face"
	} else if itemId >= 1080000 && itemId < 1090000 {
		return "Glove"
	} else if itemId >= 30000 && itemId < 35000 {
		return "Hair"
	} else if itemId >= 1050000 && itemId < 1060000 {
		return "Longcoat"
	} else if itemId >= 1060000 && itemId < 1070000 {
		return "Pants"
	} else if itemId >= 1610000 && itemId < 1660000 {
		return "Mechanic"
	} else if itemId >= 1802000 && itemId < 1820000 {
		return "PetEquip"
	} else if itemId >= 1920000 && itemId < 2000000 {
		return "Dragon"
	} else if itemId >= 1112000 && itemId < 1120000 {
		return "Ring"
	} else if itemId >= 1092000 && itemId < 1100000 {
		return "Shield"
	} else if itemId >= 1070000 && itemId < 1080000 {
		return "Shoes"
	} else if itemId >= 1900000 && itemId < 1920000 {
		return "Taming"
	} else if itemId >= 1300000 && itemId < 1800000 {
		return "Weapon"
	}
	return "Weapon"
}
