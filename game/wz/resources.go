// Package data provides MapleStory game data specifications and types.
// This file contains resource loading and management functionality.
package wz

import (
	"encoding/xml"
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// WZNode represents a node in the WZ XML structure
type WZNode struct {
	XMLName xml.Name
	Name    string     `xml:"name,attr"`
	Value   string     `xml:"value,attr,omitempty"`
	Nodes   []WZNode   `xml:"imgdir,omitempty"`
	Strings []WZString `xml:"string,omitempty"`
	Ints    []WZInt    `xml:"int,omitempty"`
	Vectors []WZVector `xml:"vector,omitempty"`
}

// WZString represents a string value in WZ structure
type WZString struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// WZInt represents an int value in WZ structure
type WZInt struct {
	Name  string `xml:"name,attr"`
	Value int    `xml:"value,attr"`
}

// WZVector represents a vector/point value in WZ structure
type WZVector struct {
	Name string `xml:"name,attr"`
	X    int    `xml:"x,attr"`
	Y    int    `xml:"y,attr"`
}

// node represents an XML node from MapleStory WZ files (legacy compatibility).
// This is kept for backward compatibility with existing code.
type node struct {
	Name     string        `xml:"name,attr"`
	Value    string        `xml:"value,attr"`
	Children []node        `xml:"imgdir"`
	Strings  []stringField `xml:"string"`
	Ints     []intField    `xml:"int"`
}

// intField represents an int element in WZ XML
type intField struct {
	Name  string `xml:"name,attr"`
	Value int    `xml:"value,attr"`
}

// stringField represents a string element in WZ XML
type stringField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// StringData contains string data for different types
type StringData struct {
	// Map strings: region -> mapId -> {mapName, streetName}
	MapStrings map[string]map[uint32]map[string]string // region -> mapId -> {mapName, streetName}

	// Mob strings: mobId -> {name}
	MobStrings map[uint32]map[string]string // mobId -> {name}

	// NPC strings: npcId -> {name}
	NpcStrings map[uint32]map[string]string // npcId -> {name}

	// Skill strings: skillId -> {name}
	SkillStrings map[uint32]map[string]string // skillId -> {name}

	// Item strings by category
	ItemCashStrings    map[uint32]map[string]string            // itemId -> {name, msg, desc}
	ItemConsumeStrings map[uint32]map[string]string            // itemId -> {name, msg, desc}
	ItemEqpStrings     map[string]map[uint32]map[string]string // category -> itemId -> {name, msg, desc}
	ItemEtcStrings     map[uint32]map[string]string            // itemId -> {name, msg, desc}
	ItemInsStrings     map[uint32]map[string]string            // itemId -> {name, msg, desc}
	ItemPetStrings     map[uint32]map[string]string            // itemId -> {name, msg, desc}
}

// Resources contains all loaded MapleStory game data.
type Resources struct {
	// Name lookup indexes (following renewal branch pattern)
	mapNameToId  map[string]uint32 // normalized map name -> map ID
	mobNameToId  map[string]uint32 // normalized mob name -> mob ID
	npcNameToId  map[string]uint32 // normalized NPC name -> NPC ID
	itemNameToId map[string]uint32 // normalized item name -> item ID

	Maps     map[uint32]*Map   // All map specifications
	Monsters map[uint32]*Mob   // All monster specifications
	Items    map[uint32]Item   // All item specifications
	Drops    map[uint32][]Drop // Monster drop tables
	Skills   map[uint32]*Skill // All skill specifications
	Strings  *StringData       // String data organized by type
	ExpTable []uint32          // Experience table: index = level, value = exp needed for that level
	Shops    map[uint32]*Shop  // NPC shops: npcId -> shop
}

// find searches for a child node by path (supports ":" separated paths).
func (node *node) find(name string) *node {

	parts := strings.Split(name, ":")
	current := node
	for _, part := range parts {
		found := false
		for _, v := range current.Children {
			if v.Name == part {
				current = &v
				found = true
				break
			}
		}

		if !found {
			return nil
		}
	}

	return current
}

// loadResourceFiles loads multiple XML files concurrently using worker goroutines.
// Generic function that processes files and calls callback with progress updates.
func loadResourceFiles[T any](root string, workerCount int, action func(path string) (result *T, err error), callback func(percent float32, value *T)) error {
	var allFiles []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(d.Name(), ".img.xml") {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	total := len(allFiles)
	if total == 0 {
		return nil
	}

	jobs := make(chan string, total)
	results := make(chan *T, total)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				m, err := action(path)

				if err == nil {
					results <- m
				} else {
					log.Println(err)
				}
			}
		}()
	}

	for _, path := range allFiles {
		jobs <- path
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	count := 0
	for m := range results {
		count++
		percent := float32(count) / float32(total) * 100
		callback(percent, m)
	}

	return nil
}

// normalizeName normalizes a name for lookup by:
// 1. Converting to lowercase (case-insensitive)
// 2. Removing all whitespace characters
func normalizeName(name string) string {
	// Convert to lowercase
	normalized := strings.ToLower(name)

	// Remove all whitespace characters
	var builder strings.Builder
	for _, r := range normalized {
		if !unicode.IsSpace(r) {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// NameToItem returns item ID by item name (case-insensitive, whitespace-insensitive)
func (r *Resources) NameToItem(name string) (uint32, bool) {
	key := normalizeName(name)
	itemId, ok := r.itemNameToId[key]
	return itemId, ok
}

// NameToMob returns mob ID by mob name (case-insensitive, whitespace-insensitive)
func (r *Resources) NameToMob(name string) (uint32, bool) {
	key := normalizeName(name)
	mobId, ok := r.mobNameToId[key]
	return mobId, ok
}

// NameToMap returns map ID by map name (case-insensitive, whitespace-insensitive)
func (r *Resources) NameToMap(name string) (uint32, bool) {
	key := normalizeName(name)
	mapId, ok := r.mapNameToId[key]
	return mapId, ok
}

// NameToNpc returns NPC ID by NPC name (case-insensitive, whitespace-insensitive)
func (r *Resources) NameToNpc(name string) (uint32, bool) {
	key := normalizeName(name)
	npcId, ok := r.npcNameToId[key]
	return npcId, ok
}

// FindWzPath attempts to find the WZ files directory by checking multiple possible paths.
// This is useful in debugging environments where the working directory may differ.
func FindWzPath(configPath string) string {
	// Try the config path as-is first
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	// Try relative to current working directory
	if absPath, err := filepath.Abs(configPath); err == nil {
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	// Try relative to executable directory
	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		candidate := filepath.Join(execDir, configPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	// Try common relative paths
	candidates := []string{
		filepath.Join("..", configPath),
		filepath.Join("..", "..", configPath),
		filepath.Join(".", configPath),
	}

	for _, candidate := range candidates {
		if absPath, err := filepath.Abs(candidate); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}

	// If nothing found, return the original path
	// (caller will handle the error)
	return configPath
}

// NewResources loads all MapleStory game data from WZ files.
// wzPath is the base path to the WZ files directory.
// Returns a fully populated Resources struct with maps, monsters, items, and drops.
func NewResources(wzPath string) *Resources {
	// Find the actual WZ path (handles debugging environments)
	originalPath := wzPath
	wzPath = FindWzPath(wzPath)
	if wzPath != originalPath {
		log.Printf("WZ path resolved: %s -> %s", originalPath, wzPath)
	} else {
		log.Printf("Loading WZ files from: %s", wzPath)
	}

	workerCount := runtime.NumCPU() * 2

	drop := map[uint32][]Drop{}
	// Load Reward.img.xml (single file, not a directory)
	rewardPath := filepath.Join(wzPath, "Reward.img.xml")
	if dropData, err := loadDrops(rewardPath); err == nil && dropData != nil {
		// Merge drop data from the file
		for mobID, drops := range *dropData {
			drop[mobID] = drops
		}
		fmt.Println("Drop files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Reward.img.xml: %v", err)
	}

	shops := map[uint32]*Shop{}
	// Load NpcShop.img.xml (single file, not a directory)
	shopPath := filepath.Join(wzPath, "NpcShop.img.xml")
	if shopData, err := loadNpcShops(shopPath); err == nil && shopData != nil {
		for npcID, shop := range *shopData {
			shops[npcID] = shop
		}
		fmt.Println("NPC shop files loaded.")
	} else if err != nil {
		log.Printf("Failed to load NpcShop.img.xml: %v", err)
	}

	items := map[uint32]Item{}
	var err error
	err = loadResourceFiles(filepath.Join(wzPath, "Character.wz"),
		workerCount,
		func(path string) (result *Equipment, err error) {
			return loadWeapons(path)
		},
		func(percent float32, value *Equipment) {
			if value != nil {
				items[value.ID] = value
			}
			fmt.Printf("Loading equipment files: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	mobs := map[uint32]*Mob{}
	err = loadResourceFiles(filepath.Join(wzPath, "Mob.wz"), workerCount, func(path string) (result *Mob, err error) {
		return loadMob(path)
	}, func(percent float32, value *Mob) {
		mobs[value.ID] = value
		fmt.Printf("Loading mob files: %.1f%%\n", percent)
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles(filepath.Join(wzPath, "Item.wz", "Consume"),
		workerCount,
		func(path string) (result *[]*Consume, err error) {
			return loadConsumes(path)
		},
		func(percent float32, value *[]*Consume) {

			for _, v := range *value {
				items[v.ID] = v
			}
			fmt.Printf("Loading consume item files: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles(filepath.Join(wzPath, "Item.wz", "Cash"),
		workerCount,
		func(path string) (result *[]*CashItem, err error) {
			return loadCashItems(path)
		},
		func(percent float32, value *[]*CashItem) {

			for _, v := range *value {
				items[v.ID] = v
			}
			fmt.Printf("Loading cash item files: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles(filepath.Join(wzPath, "Item.wz", "Install"),
		workerCount,
		func(path string) (result *[]*Installation, err error) {
			return loadInstallations(path)
		},
		func(percent float32, value *[]*Installation) {

			for _, v := range *value {
				items[v.ID] = v
			}
			fmt.Printf("Loading installation item files: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles(filepath.Join(wzPath, "Item.wz", "Special"),
		workerCount,
		func(path string) (result *[]*SpecialItem, err error) {
			return loadSpecialItems(path)
		},
		func(percent float32, value *[]*SpecialItem) {

			for _, v := range *value {
				items[v.ID] = v
			}
			fmt.Printf("Loading special item files: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles(filepath.Join(wzPath, "Item.wz", "Etc"),
		workerCount,
		func(path string) (result *[]*GeneralItem, err error) {
			return loadGeneralItems(path)
		},
		func(percent float32, value *[]*GeneralItem) {

			for _, v := range *value {
				items[v.ID] = v
			}
			fmt.Printf("Loading general item files: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = loadResourceFiles(filepath.Join(wzPath, "Item.wz", "Pet"),
		workerCount,
		func(path string) (result *Pet, err error) {
			return loadPets(path)
		},
		func(percent float32, value *Pet) {

			items[value.ID] = value
			fmt.Printf("Loading pet files: %.1f%%\n", percent)
		})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Load string data organized by image file type
	stringData := &StringData{
		MapStrings:         make(map[string]map[uint32]map[string]string),
		MobStrings:         make(map[uint32]map[string]string),
		NpcStrings:         make(map[uint32]map[string]string),
		SkillStrings:       make(map[uint32]map[string]string),
		ItemCashStrings:    make(map[uint32]map[string]string),
		ItemConsumeStrings: make(map[uint32]map[string]string),
		ItemEqpStrings:     make(map[string]map[uint32]map[string]string),
		ItemEtcStrings:     make(map[uint32]map[string]string),
		ItemInsStrings:     make(map[uint32]map[string]string),
		ItemPetStrings:     make(map[uint32]map[string]string),
	}

	// Load Map.img.xml (single file, not a directory)
	mapPath := filepath.Join(wzPath, "String.wz", "Map.img.xml")
	if mapData, err := loadMapStrings(mapPath); err == nil && mapData != nil {
		for region, regionData := range *mapData {
			stringData.MapStrings[region] = regionData
		}
		fmt.Println("Map string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Map.img.xml: %v", err)
	}

	// Load Mob.img.xml (single file, not a directory)
	mobPath := filepath.Join(wzPath, "String.wz", "Mob.img.xml")
	if mobData, err := loadMobStrings(mobPath); err == nil && mobData != nil {
		maps.Copy(stringData.MobStrings, *mobData)
		fmt.Println("Mob string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Mob.img.xml: %v", err)
	}

	// Load Npc.img.xml (single file, not a directory)
	npcPath := filepath.Join(wzPath, "String.wz", "Npc.img.xml")
	if npcData, err := loadNpcStrings(npcPath); err == nil && npcData != nil {
		maps.Copy(stringData.NpcStrings, *npcData)
		fmt.Println("NPC string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Npc.img.xml: %v", err)
	}

	// Load Skill.img.xml (single file, not a directory)
	skillPath := filepath.Join(wzPath, "String.wz", "Skill.img.xml")
	if skillData, err := loadSkillStrings(skillPath); err == nil && skillData != nil {
		maps.Copy(stringData.SkillStrings, *skillData)
		fmt.Println("Skill string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Skill.img.xml: %v", err)
	}

	// Load Cash.img.xml (single file, not a directory)
	cashPath := filepath.Join(wzPath, "String.wz", "Cash.img.xml")
	if cashData, err := loadItemStrings(cashPath); err == nil && cashData != nil {
		maps.Copy(stringData.ItemCashStrings, *cashData)
		fmt.Println("Cash item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Cash.img.xml: %v", err)
	}

	// Load Consume.img.xml (single file, not a directory)
	consumePath := filepath.Join(wzPath, "String.wz", "Consume.img.xml")
	if consumeData, err := loadItemStrings(consumePath); err == nil && consumeData != nil {
		maps.Copy(stringData.ItemConsumeStrings, *consumeData)
		fmt.Println("Consume item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Consume.img.xml: %v", err)
	}

	// Load Eqp.img.xml (single file, not a directory)
	eqpPath := filepath.Join(wzPath, "String.wz", "Eqp.img.xml")
	if eqpData, err := loadEqpStrings(eqpPath); err == nil && eqpData != nil {
		for category, categoryData := range *eqpData {
			stringData.ItemEqpStrings[category] = categoryData
		}
		fmt.Println("Equipment item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Eqp.img.xml: %v", err)
	}

	// Load Etc.img.xml (single file, not a directory)
	etcPath := filepath.Join(wzPath, "String.wz", "Etc.img.xml")
	if etcData, err := loadItemStrings(etcPath); err == nil && etcData != nil {
		maps.Copy(stringData.ItemEtcStrings, *etcData)
		fmt.Println("Etc item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Etc.img.xml: %v", err)
	}

	// Load Ins.img.xml (single file, not a directory)
	insPath := filepath.Join(wzPath, "String.wz", "Ins.img.xml")
	if insData, err := loadItemStrings(insPath); err == nil && insData != nil {
		maps.Copy(stringData.ItemInsStrings, *insData)
		fmt.Println("Installation item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Ins.img.xml: %v", err)
	}

	// Load Pet.img.xml (single file, not a directory)
	petPath := filepath.Join(wzPath, "String.wz", "Pet.img.xml")
	if petData, err := loadItemStrings(petPath); err == nil && petData != nil {
		maps.Copy(stringData.ItemPetStrings, *petData)
		fmt.Println("Pet item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Pet.img.xml: %v", err)
	}

	maps := map[uint32]*Map{}
	err = loadResourceFiles(filepath.Join(wzPath, "Map.wz", "Map"), workerCount, func(path string) (result *Map, err error) {
		base := filepath.Base(path)
		idStr := strings.TrimSuffix(base, ".img.xml")
		mapId, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}

		m, err := loadMaps(path, uint32(mapId))
		if err != nil {
			return nil, err
		}
		return m, nil
	}, func(percent float32, value *Map) {
		maps[value.ID] = value
		fmt.Printf("Loading map files: %.1f%%\n", percent)
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	skills := map[uint32]*Skill{}
	skillPath = filepath.Join(wzPath, "Skill.wz")
	var skillFiles []string
	err = filepath.WalkDir(skillPath, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".img.xml") {
			base := d.Name()
			// Skip special files like ItemSkill.img.xml, MobSkill.img.xml, MCSkill.img.xml
			if base != "ItemSkill.img.xml" && base != "MobSkill.img.xml" && base != "MCSkill.img.xml" && base != "MCGuardian.img.xml" {
				skillFiles = append(skillFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	total := len(skillFiles)
	for i, path := range skillFiles {
		jobSkills, loadErr := loadSkillJobFile(path)
		if loadErr != nil {
			log.Printf("Failed to load skill file %s: %v", path, loadErr)
			continue
		}
		for skillID, skill := range jobSkills {
			skills[skillID] = skill
		}
		if total > 0 {
			percent := float32(i+1) * 100.0 / float32(total)
			fmt.Printf("Loading skill files: %.1f%%\n", percent)
		}
	}

	fmt.Println("\nAll files loaded.")

	expTable := getHardcodedExpTable()

	result := &Resources{
		mapNameToId:  make(map[string]uint32),
		mobNameToId:  make(map[string]uint32),
		npcNameToId:  make(map[string]uint32),
		itemNameToId: make(map[string]uint32),
		Maps:         maps,
		Monsters:     mobs,
		Items:        items,
		Drops:        drop,
		Strings:      stringData,
		ExpTable:     expTable,
		Skills:       skills,
		Shops:        shops,
	}

	// Build name lookup indexes from string data (following renewal branch pattern)
	result.buildNameIndexes()

	return result
}

// buildNameIndexes builds all name to ID indexes (following renewal branch pattern)
func (r *Resources) buildNameIndexes() {
	r.buildMapNameIndex()
	r.buildMobNameIndex()
	r.buildNpcNameIndex()
	r.buildItemNameIndex()
}

// GetSkill returns a skill by ID
func (r *Resources) GetSkill(skillID uint32) *Skill {
	return r.Skills[skillID]
}

// GetShop returns a shop by NPC ID
func (r *Resources) GetShop(npcID uint32) *Shop {
	return r.Shops[npcID]
}

// GetExpNeededForLevel returns the cumulative experience needed to reach the specified level.
// expTable[level] contains the cumulative exp needed to reach level (level+1).
func (r *Resources) GetExpNeededForLevel(level uint8) uint32 {
	if level <= 0 || level > 200 {
		return 0
	}
	if int(level) >= len(r.ExpTable) {
		return 0
	}
	return r.ExpTable[level]
}

// buildMapNameIndex builds the name to ID index for maps
// Only includes maps that are actually loaded in r.Maps (from Map.wz)
// When multiple maps have the same name, the smallest map ID is prioritized
func (r *Resources) buildMapNameIndex() {
	// Collect all map IDs and sort them to ensure deterministic behavior
	mapIds := make([]uint32, 0, len(r.Maps))
	for mapId := range r.Maps {
		mapIds = append(mapIds, mapId)
	}

	// Sort map IDs to ensure smaller IDs are processed first
	// This ensures consistent behavior when multiple maps have the same name
	sort.Slice(mapIds, func(i, j int) bool {
		return mapIds[i] < mapIds[j]
	})

	// Iterate through sorted map IDs
	for _, mapId := range mapIds {
		// Get map name from String.wz
		for _, regionMaps := range r.Strings.MapStrings {
			if mapNameData, ok := regionMaps[mapId]; ok && mapNameData != nil {
				if mapName, ok := mapNameData["mapName"]; ok && mapName != "" {
					key := normalizeName(mapName)
					if _, exists := r.mapNameToId[key]; !exists {
						r.mapNameToId[key] = mapId
					}
				}
				break // Found the map, no need to check other regions
			}
		}
	}
}

// buildMobNameIndex builds the name to ID index for mobs
// Only includes mobs that are actually loaded in r.Monsters (from Mob.wz)
// Prioritizes original mobs (without link) over linked mobs
func (r *Resources) buildMobNameIndex() {
	// First pass: Add original mobs (without link) to index
	for mobId, mob := range r.Monsters {
		// Skip linked mobs in first pass
		if mob.Link != "" {
			continue
		}
		// Get mob name from String.wz
		if mobNameData, ok := r.Strings.MobStrings[mobId]; ok && mobNameData != nil {
			if mobName, ok := mobNameData["name"]; ok && mobName != "" {
				// Normalize name (lowercase + remove whitespace)
				key := normalizeName(mobName)
				// Add to index (original mobs have priority)
				if _, exists := r.mobNameToId[key]; !exists {
					r.mobNameToId[key] = mobId
				}
			}
		}
	}

	// Second pass: Add linked mobs to index (only if name not already exists)
	for mobId, mob := range r.Monsters {
		// Only process linked mobs in second pass
		if mob.Link == "" {
			continue
		}
		// Get mob name from String.wz
		if mobNameData, ok := r.Strings.MobStrings[mobId]; ok && mobNameData != nil {
			if mobName, ok := mobNameData["name"]; ok && mobName != "" {
				// Normalize name (lowercase + remove whitespace)
				key := normalizeName(mobName)
				// Only add if name not already exists (original mobs have priority)
				if _, exists := r.mobNameToId[key]; !exists {
					r.mobNameToId[key] = mobId
				}
			}
		}
	}
}

// buildNpcNameIndex builds the name to ID index for NPCs
// Only includes NPCs that are actually spawned in loaded maps (from Map.wz)
func (r *Resources) buildNpcNameIndex() {
	// Collect all NPC IDs from loaded maps
	npcIds := make(map[uint32]bool)
	for _, mapData := range r.Maps {
		for _, spawn := range mapData.NpcSpawns {
			npcIds[spawn.ID] = true
		}
	}

	// Iterate through NPCs that are actually spawned in maps
	for npcId := range npcIds {
		// Get NPC name from String.wz
		if npcNameData, ok := r.Strings.NpcStrings[npcId]; ok && npcNameData != nil {
			if npcName, ok := npcNameData["name"]; ok && npcName != "" {
				// Normalize name (lowercase + remove whitespace)
				key := normalizeName(npcName)
				// If multiple NPCs have the same name, keep the first one found
				if _, exists := r.npcNameToId[key]; !exists {
					r.npcNameToId[key] = npcId
				}
			}
		}
	}
}

// buildItemNameIndex builds the name to ID index for items
// Only includes items that are actually loaded in r.Items (from Item.wz)
func (r *Resources) buildItemNameIndex() {
	// Iterate through actually loaded items (from Item.wz)
	for itemId := range r.Items {
		// Get item name from String.wz based on item ID range (same logic as GetItemName)
		itemName := ""

		if itemId >= 5010000 {
			// Cash items
			if itemNameData, ok := r.Strings.ItemCashStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else if itemId >= 2000000 && itemId < 3000000 {
			// Consumable items
			if itemNameData, ok := r.Strings.ItemConsumeStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else if itemId >= 4000000 && itemId < 5000000 {
			// Etc items
			if itemNameData, ok := r.Strings.ItemEtcStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else if itemId >= 3000000 && itemId < 4000000 {
			// Installation items
			if itemNameData, ok := r.Strings.ItemInsStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else if itemId >= 5000000 && itemId < 5010000 {
			// Pet items
			if itemNameData, ok := r.Strings.ItemPetStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else {
			// Equipment items - need to find category using getItemCategory
			category := getItemCategory(itemId)
			if categoryData, ok := r.Strings.ItemEqpStrings[category]; ok {
				if itemNameData, ok := categoryData[itemId]; ok && itemNameData != nil {
					itemName, _ = itemNameData["name"]
				}
			}
		}

		if itemName != "" {
			key := normalizeName(itemName)
			if _, exists := r.itemNameToId[key]; !exists {
				r.itemNameToId[key] = itemId
			}
		}
	}
}
