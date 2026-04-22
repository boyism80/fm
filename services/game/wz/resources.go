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

type WZNode struct {
	XMLName xml.Name
	Name    string     `xml:"name,attr"`
	Value   string     `xml:"value,attr,omitempty"`
	Nodes   []WZNode   `xml:"imgdir,omitempty"`
	Strings []WZString `xml:"string,omitempty"`
	Ints    []WZInt    `xml:"int,omitempty"`
	Vectors []WZVector `xml:"vector,omitempty"`
}

type WZString struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type WZInt struct {
	Name  string `xml:"name,attr"`
	Value int    `xml:"value,attr"`
}

type WZVector struct {
	Name string `xml:"name,attr"`
	X    int    `xml:"x,attr"`
	Y    int    `xml:"y,attr"`
}

type node struct {
	Name     string        `xml:"name,attr"`
	Value    string        `xml:"value,attr"`
	Children []node        `xml:"imgdir"`
	Strings  []stringField `xml:"string"`
	Ints     []intField    `xml:"int"`
	Floats   []floatField  `xml:"float"`
	Vectors  []vectorField `xml:"vector"`
}

type floatField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type intField struct {
	Name  string `xml:"name,attr"`
	Value int    `xml:"value,attr"`
}

type stringField struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type vectorField struct {
	Name string `xml:"name,attr"`
	X    int    `xml:"x,attr"`
	Y    int    `xml:"y,attr"`
}

type StringData struct {
	MapStrings map[string]map[uint32]map[string]string

	MobStrings map[uint32]map[string]string

	NpcStrings map[uint32]map[string]string

	SkillStrings map[uint32]map[string]string

	ItemCashStrings    map[uint32]map[string]string
	ItemConsumeStrings map[uint32]map[string]string
	ItemEqpStrings     map[string]map[uint32]map[string]string
	ItemEtcStrings     map[uint32]map[string]string
	ItemInsStrings     map[uint32]map[string]string
	ItemPetStrings     map[uint32]map[string]string
}

type Resources struct {
	mapNameToId   map[string]uint32
	mobNameToId   map[string]uint32
	npcNameToId   map[string]uint32
	itemNameToId  map[string]uint32
	skillNameToId map[string]uint32

	Maps     map[uint32]*Map
	Monsters map[uint32]*Mob
	Items    map[uint32]Item
	Drops    map[uint32][]Drop
	Skills   map[uint32]*Skill
	Strings  *StringData
	ExpTable []uint32
	Shops    map[uint32]*Shop
}

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

func loadEquipmentFiles(root string, workerCount int, items map[uint32]Item) error {
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
	classes := make(chan string, total)
	results := make(chan Item, total)
	var wg sync.WaitGroup
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range classes {
				item, loadErr := loadWeapons(path)
				if loadErr == nil && item != nil {
					results <- item
				} else if loadErr != nil {
					log.Println(loadErr)
				}
			}
		}()
	}
	go func() {
		for _, path := range allFiles {
			classes <- path
		}
		close(classes)
		wg.Wait()
		close(results)
	}()
	count := 0
	for item := range results {
		count++
		items[item.GetID()] = item
		if total > 0 && count%500 == 0 {
			fmt.Printf("Loading equipment files: %.1f%%\n", float32(count)/float32(total)*100)
		}
	}
	fmt.Printf("Loading equipment files: 100.0%%\n")
	return nil
}

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

	classes := make(chan string, total)
	results := make(chan *T, total)
	var wg sync.WaitGroup

	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range classes {
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
		classes <- path
	}
	close(classes)

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

func normalizeName(name string) string {

	normalized := strings.ToLower(name)

	var builder strings.Builder
	for _, r := range normalized {
		if !unicode.IsSpace(r) {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

func (r *Resources) NameToItem(name string) (uint32, bool) {
	key := normalizeName(name)
	itemId, ok := r.itemNameToId[key]
	return itemId, ok
}

func (r *Resources) NameToMob(name string) (uint32, bool) {
	key := normalizeName(name)
	mobId, ok := r.mobNameToId[key]
	return mobId, ok
}

func (r *Resources) NameToMap(name string) (uint32, bool) {
	key := normalizeName(name)
	mapId, ok := r.mapNameToId[key]
	return mapId, ok
}

func (r *Resources) NameToNpc(name string) (uint32, bool) {
	key := normalizeName(name)
	npcId, ok := r.npcNameToId[key]
	return npcId, ok
}

func (r *Resources) NameToSkill(name string) (uint32, bool) {
	key := normalizeName(name)
	skillID, ok := r.skillNameToId[key]
	return skillID, ok
}

func FindWzPath(configPath string) string {

	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	if absPath, err := filepath.Abs(configPath); err == nil {
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	if execPath, err := os.Executable(); err == nil {
		execDir := filepath.Dir(execPath)
		candidate := filepath.Join(execDir, configPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

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

	return configPath
}

func NewResources(wzPath string) *Resources {

	originalPath := wzPath
	wzPath = FindWzPath(wzPath)
	if wzPath != originalPath {
		log.Printf("WZ path resolved: %s -> %s", originalPath, wzPath)
	} else {
		log.Printf("Loading WZ files from: %s", wzPath)
	}

	workerCount := runtime.NumCPU() * 2

	drop := map[uint32][]Drop{}

	rewardPath := filepath.Join(wzPath, "Reward.img.xml")
	if dropData, err := loadDrops(rewardPath); err == nil && dropData != nil {

		for mobID, drops := range *dropData {
			drop[mobID] = drops
		}
		fmt.Println("Drop files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Reward.img.xml: %v", err)
	}

	shops := map[uint32]*Shop{}

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
	err = loadEquipmentFiles(filepath.Join(wzPath, "Character.wz"), workerCount, items)
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
		func(path string) (result *[]*MiscItem, err error) {
			return loadMiscItems(path)
		},
		func(percent float32, value *[]*MiscItem) {

			for _, v := range *value {
				items[v.ID] = v
			}
			fmt.Printf("Loading misc item files: %.1f%%\n", percent)
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

	mapPath := filepath.Join(wzPath, "String.wz", "Map.img.xml")
	if mapData, err := loadMapStrings(mapPath); err == nil && mapData != nil {
		for region, regionData := range *mapData {
			stringData.MapStrings[region] = regionData
		}
		fmt.Println("Map string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Map.img.xml: %v", err)
	}

	mobPath := filepath.Join(wzPath, "String.wz", "Mob.img.xml")
	if mobData, err := loadMobStrings(mobPath); err == nil && mobData != nil {
		maps.Copy(stringData.MobStrings, *mobData)
		fmt.Println("Mob string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Mob.img.xml: %v", err)
	}

	npcPath := filepath.Join(wzPath, "String.wz", "Npc.img.xml")
	if npcData, err := loadNpcStrings(npcPath); err == nil && npcData != nil {
		maps.Copy(stringData.NpcStrings, *npcData)
		fmt.Println("NPC string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Npc.img.xml: %v", err)
	}

	skillPath := filepath.Join(wzPath, "String.wz", "Skill.img.xml")
	if skillData, err := loadSkillStrings(skillPath); err == nil && skillData != nil {
		maps.Copy(stringData.SkillStrings, *skillData)
		fmt.Println("Skill string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Skill.img.xml: %v", err)
	}

	cashPath := filepath.Join(wzPath, "String.wz", "Cash.img.xml")
	if cashData, err := loadItemStrings(cashPath); err == nil && cashData != nil {
		maps.Copy(stringData.ItemCashStrings, *cashData)
		fmt.Println("Cash item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Cash.img.xml: %v", err)
	}

	consumePath := filepath.Join(wzPath, "String.wz", "Consume.img.xml")
	if consumeData, err := loadItemStrings(consumePath); err == nil && consumeData != nil {
		maps.Copy(stringData.ItemConsumeStrings, *consumeData)
		fmt.Println("Consume item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Consume.img.xml: %v", err)
	}

	eqpPath := filepath.Join(wzPath, "String.wz", "Eqp.img.xml")
	if eqpData, err := loadEqpStrings(eqpPath); err == nil && eqpData != nil {
		for category, categoryData := range *eqpData {
			stringData.ItemEqpStrings[category] = categoryData
		}
		fmt.Println("Equipment item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Eqp.img.xml: %v", err)
	}

	etcPath := filepath.Join(wzPath, "String.wz", "Etc.img.xml")
	if etcData, err := loadItemStrings(etcPath); err == nil && etcData != nil {
		maps.Copy(stringData.ItemEtcStrings, *etcData)
		fmt.Println("Etc item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Etc.img.xml: %v", err)
	}

	insPath := filepath.Join(wzPath, "String.wz", "Ins.img.xml")
	if insData, err := loadItemStrings(insPath); err == nil && insData != nil {
		maps.Copy(stringData.ItemInsStrings, *insData)
		fmt.Println("Installation item string files loaded.")
	} else if err != nil {
		log.Printf("Failed to load Ins.img.xml: %v", err)
	}

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
		classSkills, loadErr := loadSkillClassFile(path)
		if loadErr != nil {
			log.Printf("Failed to load skill file %s: %v", path, loadErr)
			continue
		}
		for skillID, skill := range classSkills {
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
		mapNameToId:   make(map[string]uint32),
		mobNameToId:   make(map[string]uint32),
		npcNameToId:   make(map[string]uint32),
		itemNameToId:  make(map[string]uint32),
		skillNameToId: make(map[string]uint32),
		Maps:          maps,
		Monsters:      mobs,
		Items:         items,
		Drops:         drop,
		Strings:       stringData,
		ExpTable:      expTable,
		Skills:        skills,
		Shops:         shops,
	}

	result.buildNameIndexes()

	return result
}

func (r *Resources) buildNameIndexes() {
	r.buildMapNameIndex()
	r.buildMobNameIndex()
	r.buildNpcNameIndex()
	r.buildItemNameIndex()
	r.buildSkillNameIndex()
}

func (r *Resources) GetSkill(skillID uint32) *Skill {
	return r.Skills[skillID]
}

func (r *Resources) GetShop(npcID uint32) *Shop {
	return r.Shops[npcID]
}

func (r *Resources) GetExpNeededForLevel(level uint8) uint32 {
	if level <= 0 || level > 200 {
		return 0
	}
	if int(level) >= len(r.ExpTable) {
		return 0
	}
	return r.ExpTable[level]
}

func (r *Resources) buildMapNameIndex() {

	mapIds := make([]uint32, 0, len(r.Maps))
	for mapId := range r.Maps {
		mapIds = append(mapIds, mapId)
	}

	sort.Slice(mapIds, func(i, j int) bool {
		return mapIds[i] < mapIds[j]
	})

	for _, mapId := range mapIds {

		for _, regionMaps := range r.Strings.MapStrings {
			if mapNameData, ok := regionMaps[mapId]; ok && mapNameData != nil {
				if mapName, ok := mapNameData["mapName"]; ok && mapName != "" {
					key := normalizeName(mapName)
					if _, exists := r.mapNameToId[key]; !exists {
						r.mapNameToId[key] = mapId
					}
				}
				break
			}
		}
	}
}

func (r *Resources) buildMobNameIndex() {

	for mobId, mob := range r.Monsters {

		if mob.Link != "" {
			continue
		}

		if mobNameData, ok := r.Strings.MobStrings[mobId]; ok && mobNameData != nil {
			if mobName, ok := mobNameData["name"]; ok && mobName != "" {

				key := normalizeName(mobName)

				if _, exists := r.mobNameToId[key]; !exists {
					r.mobNameToId[key] = mobId
				}
			}
		}
	}

	for mobId, mob := range r.Monsters {

		if mob.Link == "" {
			continue
		}

		if mobNameData, ok := r.Strings.MobStrings[mobId]; ok && mobNameData != nil {
			if mobName, ok := mobNameData["name"]; ok && mobName != "" {

				key := normalizeName(mobName)

				if _, exists := r.mobNameToId[key]; !exists {
					r.mobNameToId[key] = mobId
				}
			}
		}
	}
}

func (r *Resources) buildNpcNameIndex() {

	npcIds := make(map[uint32]bool)
	for _, mapData := range r.Maps {
		for _, spawn := range mapData.NpcSpawns {
			npcIds[spawn.ID] = true
		}
	}

	for npcId := range npcIds {

		if npcNameData, ok := r.Strings.NpcStrings[npcId]; ok && npcNameData != nil {
			if npcName, ok := npcNameData["name"]; ok && npcName != "" {

				key := normalizeName(npcName)

				if _, exists := r.npcNameToId[key]; !exists {
					r.npcNameToId[key] = npcId
				}
			}
		}
	}
}

func (r *Resources) buildItemNameIndex() {

	for itemId := range r.Items {

		itemName := ""

		if itemId >= 5010000 {

			if itemNameData, ok := r.Strings.ItemCashStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else if itemId >= 2000000 && itemId < 3000000 {

			if itemNameData, ok := r.Strings.ItemConsumeStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else if itemId >= 4000000 && itemId < 5000000 {

			if itemNameData, ok := r.Strings.ItemEtcStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else if itemId >= 3000000 && itemId < 4000000 {

			if itemNameData, ok := r.Strings.ItemInsStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else if itemId >= 5000000 && itemId < 5010000 {

			if itemNameData, ok := r.Strings.ItemPetStrings[itemId]; ok && itemNameData != nil {
				itemName, _ = itemNameData["name"]
			}
		} else {

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

func (r *Resources) buildSkillNameIndex() {
	for skillID := range r.Skills {
		skillNameData, ok := r.Strings.SkillStrings[skillID]
		if !ok || skillNameData == nil {
			continue
		}
		skillName, ok := skillNameData["name"]
		if !ok || skillName == "" {
			continue
		}
		key := normalizeName(skillName)
		if _, exists := r.skillNameToId[key]; !exists {
			r.skillNameToId[key] = skillID
		}
	}
}
