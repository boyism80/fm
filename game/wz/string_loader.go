// Package data provides MapleStory game data specifications and types.
// This file contains string loading functions for String.wz files.
package wz

import (
	"encoding/xml"
	"os"
	"strconv"
)

// CountStrings returns the total count of all string entries
func (sd *StringData) CountStrings() int {
	count := 0

	// Count map strings
	for _, regionData := range sd.MapStrings {
		count += len(regionData)
	}

	// Count other strings
	count += len(sd.MobStrings)
	count += len(sd.NpcStrings)
	count += len(sd.SkillStrings)
	count += len(sd.ItemCashStrings)
	count += len(sd.ItemConsumeStrings)
	count += len(sd.ItemEtcStrings)
	count += len(sd.ItemInsStrings)
	count += len(sd.ItemPetStrings)

	// Count equipment strings
	for _, categoryData := range sd.ItemEqpStrings {
		count += len(categoryData)
	}

	return count
}

// loadMapStrings loads map name strings from Map.img.xml
// Structure: region -> mapId -> {mapName, streetName}
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

	// Map.img structure: region -> mapId -> {mapName, streetName}
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
			// Process string fields
			for _, field := range mapNode.Strings {
				stringData[field.Name] = field.Value
			}
			result[region][uint32(mapId)] = stringData
		}
	}

	return &result, nil
}

// loadMobStrings loads mob name strings from Mob.img.xml
// Structure: mobId -> {name}
func loadMobStrings(path string) (*map[uint32]map[string]string, error) {
	return loadSimpleStrings(path)
}

// loadNpcStrings loads NPC name strings from Npc.img.xml
// Structure: npcId -> {name}
func loadNpcStrings(path string) (*map[uint32]map[string]string, error) {
	return loadSimpleStrings(path)
}

// loadSkillStrings loads skill name strings from Skill.img.xml
// Structure: skillId (7-digit padded) -> {name}
func loadSkillStrings(path string) (*map[uint32]map[string]string, error) {
	return loadSimpleStrings(path)
}

// loadItemStrings loads item name strings from Cash.img.xml, Consume.img.xml, etc.
// Structure: itemId -> {name, msg, desc}
func loadItemStrings(path string) (*map[uint32]map[string]string, error) {
	return loadSimpleStrings(path)
}

// loadEqpStrings loads equipment name strings from Eqp.img.xml
// Structure: category -> itemId -> {name, msg, desc}
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

	// Eqp.img structure: Eqp.img -> Eqp -> category (Face, Cap, Weapon, etc.) -> itemId -> {name, msg, desc}
	// We need to skip the "Eqp" wrapper and process categories directly
	var processNode func(node *node, category string)
	processNode = func(n *node, category string) {
		// Check if this node's name is a numeric item ID
		if itemId, err := strconv.ParseUint(n.Name, 10, 32); err == nil {
			// This is an item ID node
			if result[category] == nil {
				result[category] = make(map[uint32]map[string]string)
			}

			stringData := make(map[string]string)
			// Process string fields
			for _, field := range n.Strings {
				stringData[field.Name] = field.Value
			}
			result[category][uint32(itemId)] = stringData
		} else {
			// This is a category or subcategory node
			// Process all children recursively
			for _, child := range n.Children {
				// Use current node's name as category if it's not numeric
				childCategory := category
				if category == "" {
					childCategory = n.Name
				}
				processNode(&child, childCategory)
			}
		}
	}

	// Process root children
	// Structure: Eqp.img -> Eqp -> categories (Face, Cap, Weapon, etc.)
	for _, rootChild := range root.Children {
		// Skip the "Eqp" wrapper node and process its children (actual categories)
		if rootChild.Name == "Eqp" {
			// Process categories directly (Face, Cap, Weapon, etc.)
			for _, categoryNode := range rootChild.Children {
				processNode(&categoryNode, categoryNode.Name)
			}
		} else {
			// Fallback: if structure is different, process normally
			processNode(&rootChild, rootChild.Name)
		}
	}

	return &result, nil
}

// loadSimpleStrings loads simple string structures (mob, npc, skill, item)
// Structure: id -> {name, ...}
// Some files may have nested structures (e.g., Etc -> itemId)
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

	// Recursive function to process nodes
	var processNode func(n *node)
	processNode = func(n *node) {
		// Try to parse as ID
		id, err := strconv.ParseUint(n.Name, 10, 32)
		if err != nil {
			// If not a number, it might be a category or nested structure
			// Recursively process children
			for i := range n.Children {
				processNode(&n.Children[i])
			}
			return
		}

		// This is an item ID node
		stringData := make(map[string]string)
		// Process string fields
		for _, field := range n.Strings {
			stringData[field.Name] = field.Value
		}
		result[uint32(id)] = stringData
	}

	// Process root children
	for i := range root.Children {
		processNode(&root.Children[i])
	}

	return &result, nil
}

// getMapRegion returns the region name for a given map ID
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
	// Add more region mappings as needed
	return "maple"
}

// GetMapName returns the map name and street name for a given map ID
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

// GetMobName returns the mob name for a given mob ID
func (r *Resources) GetMobName(mobId uint32) string {
	if mobData, ok := r.Strings.MobStrings[mobId]; ok {
		return mobData["name"]
	}
	return ""
}

// GetNpcName returns the NPC name for a given NPC ID
func (r *Resources) GetNpcName(npcId uint32) string {
	if npcData, ok := r.Strings.NpcStrings[npcId]; ok {
		return npcData["name"]
	}
	return ""
}

// GetSkillName returns the skill name for a given skill ID
func (r *Resources) GetSkillName(skillId uint32) string {
	if skillData, ok := r.Strings.SkillStrings[skillId]; ok {
		return skillData["name"]
	}
	return ""
}

// GetItemName returns the item name for a given item ID
func (r *Resources) GetItemName(itemId uint32) string {
	// Check based on item ID range
	if itemId >= 5010000 {
		// Cash items
		if itemData, ok := r.Strings.ItemCashStrings[itemId]; ok {
			return itemData["name"]
		}
	} else if itemId >= 2000000 && itemId < 3000000 {
		// Consume items
		if itemData, ok := r.Strings.ItemConsumeStrings[itemId]; ok {
			return itemData["name"]
		}
	} else if itemId >= 4000000 && itemId < 5000000 {
		// Etc items
		if itemData, ok := r.Strings.ItemEtcStrings[itemId]; ok {
			return itemData["name"]
		}
	} else if itemId >= 3000000 && itemId < 4000000 {
		// Installation items
		if itemData, ok := r.Strings.ItemInsStrings[itemId]; ok {
			return itemData["name"]
		}
	} else if itemId >= 5000000 && itemId < 5010000 {
		// Pet items
		if itemData, ok := r.Strings.ItemPetStrings[itemId]; ok {
			return itemData["name"]
		}
	} else {
		// Equipment items - need to find category
		category := getItemCategory(itemId)
		if categoryData, ok := r.Strings.ItemEqpStrings[category]; ok {
			if itemData, ok := categoryData[itemId]; ok {
				return itemData["name"]
			}
		}
	}
	return ""
}

// getItemCategory returns the equipment category for a given item ID
func getItemCategory(itemId uint32) string {
	// Based on String.wz.md documentation
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
	return "Weapon" // Default
}
