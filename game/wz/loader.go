// Package data provides MapleStory game data specifications and types.
// This file contains XML data loading functions for various game resources.
package wz

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/boyism80/fm/types"
)

var mutex sync.Mutex = sync.Mutex{}
var visit map[string]bool = map[string]bool{}

// loadCashItems loads cash shop item specifications from XML file.
func loadCashItems(path string) (*[]*CashItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*CashItem{}
	for _, v := range root.Children {
		model := CashItem{
			ItemCore: &ItemCore{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		model.ID = uint32(id)
		info := v.find("info")
		if info != nil {
			for _, intField := range info.Ints {
				switch intField.Name {
				case "slotMax":
					model.SlotMax = uint16(intField.Value)
				case "cash":
				}
			}

			for _, iv := range info.Children {
				switch iv.Name {
				case "icon":
				case "iconRaw":
				case "recoveryRate":
				case "npc":
				case "rate":
				case "time":
				case "path":
				case "meso":
				case "life":
				case "sample":
				case "addTime":
				case "maxDays":
				case "pickupItem":
				case "add":
				case "consumeHP":
				case "longRange":
				case "dropSweep":
				case "pickupAll":
				case "ignorePickup":
				case "consumeMP":
				case "type":
				case "floatType":
				case "stateChangeItem":
				case "direction":
				case "speed":
				case "isBgmOrEffect":
				case "bgmPath":
				case "repeat":
				case "soldInform":
				case "noFlip":
					break

				default:
					mutex.Lock()
					if _, ok := visit[iv.Name]; ok {
						mutex.Unlock()
						continue
					}
					visit[iv.Name] = true
					mutex.Unlock()
					log.Printf("%s is not declared in %s:info\n", iv.Name, filepath.Base(path))
				}
			}
		}

		specs = append(specs, &model)
	}

	return &specs, nil
}

// loadConsumes loads consumable item specifications from XML file.
func loadConsumes(path string) (*[]*Consume, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*Consume{}
	for _, v := range root.Children {
		model := Consume{
			ItemCore: &ItemCore{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		model.ID = uint32(id)
		nodeInfo := v.find("info")
		if nodeInfo != nil {
			for _, intField := range nodeInfo.Ints {
				switch intField.Name {
				case "price":
					model.Price = intField.Value
				case "slotMax":
					model.SlotMax = uint16(intField.Value)
				case "tradeBlock":
				case "notSale":
				case "only":
				case "timeLimited":
				case "quest":
				case "incPAD":
				case "reqLevel":
				case "tradeAvailable":
				case "incPDD":
				case "incMDD":
				case "incACC":
				case "incMHP":
				case "incINT":
				case "incMAD":
				case "incDEX":
				case "incLUK":
				case "incSTR":
				case "incSpeed":
				case "incMMP":
				case "incEVA":
				case "incJump":
				case "tradBlock":
				case "bigSize":
				case "scanTradeBlock":
				case "mcType":
				case "cursed":
				case "preventslip":
				case "warmsupport":
				case "reqRUC":
				case "recover":
				case "randstat":
				case "unitPrice":
				case "useDelay":
				case "delayMsg":
				case "reqSkillLevel":
				case "success":
				case "masterLevel":
				case "mobHP":
				case "bridleMsgType":
				case "bridleProp":
				case "bridlePropChg":
				case "noCancelMouse":
				}
			}

			for _, iv := range nodeInfo.Children {
				switch iv.Name {
				case "icon":
				case "iconRaw":
				case "pquest":
				case "skill":
				case "mob":
				case "create":
				case "left":
				case "right":
				case "top":
				case "bottom":
				case "type":
				case "effect":
				case "monsterBook":
					break

				default:
					mutex.Lock()
					if _, ok := visit[iv.Name]; ok {
						mutex.Unlock()
						continue
					}
					visit[iv.Name] = true
					mutex.Unlock()
					log.Printf("%s is not declared in %s:info\n", iv.Name, filepath.Base(path))
				}
			}
		}

		node := v.find("model")
		if node != nil {
			for _, sv := range node.Children {
				switch sv.Name {
				case "0":
				case "1":
				case "2":
				case "3":
				case "4":
				case "5":
				case "6":
				case "7":
				case "8":
				case "9":
					break

				case "acc":
				case "accRate":
				case "barrier":
				case "con":
				case "consumeOnPickup":
				case "cp":
				case "curse":
				case "darkness":
				case "defenseAtt":
				case "defenseState":
				case "dojangshield":
				case "eva":
				case "evaRate":
				case "expinc":
				case "ghost":
				case "hp":
				case "hpR":
				case "ignoreContinent":
				case "inc":
				case "incFatigue":
				case "itemCode":
				case "itemRange":
				case "itemupbyitem":
				case "jump":
				case "mad":
				case "madRate":
				case "mdd":
				case "mddRate":
				case "mesoupbyitem":
				case "mhpR":
				case "mhpRRate":
				case "mmpR":
				case "mmpRRate":
				case "morph":
				case "moveTo":
				case "mp":
				case "mpR":
				case "npc":
				case "nuffSkill":
				case "onlyPickup":
				case "pad":
				case "padRate":
				case "party":
				case "pdd":
				case "pddRate":
				case "poison":
				case "prob":
				case "randomMoveInFieldSet":
				case "respectFS":
				case "respectMimmune":
				case "respectPimmune":
				case "returnMapQR":
				case "script":
				case "seal":
				case "speed":
				case "speedRate":
				case "thaw":
				case "time":
				case "weakness":
				case "mob":
					break

				default:
					mutex.Lock()
					if _, ok := visit[sv.Name]; ok {
						mutex.Unlock()
						continue
					}
					visit[sv.Name] = true
					mutex.Unlock()
					log.Printf("%s is not declared in %s:model\n", sv.Name, filepath.Base(path))
				}
			}
		}

		specs = append(specs, &model)
	}

	return &specs, nil
}

// loadWeapons loads weapon and equipment specifications from XML file.
func loadWeapons(path string) (*Equipment, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	model := Equipment{
		ItemCore: &ItemCore{},
	}
	id, err := strconv.Atoi(strings.TrimSuffix(root.Name, ".img"))
	if err != nil {
		// Skip category files (hit, bow, axe, etc.) - these are not equipment items
		return nil, nil
	}
	model.ID = uint32(id)
	node := root.find("info")
	if node == nil {
		return nil, fmt.Errorf("'info' does not exist in %s", path)
	}

	for _, intField := range node.Ints {
		switch intField.Name {
		case "reqJob":
			model.Required.Class = intField.Value
		case "reqLevel":
			model.Required.Level = uint8(intField.Value)
		case "reqSTR":
			model.Required.Str = uint16(intField.Value)
		case "reqDEX":
			model.Required.Dex = uint16(intField.Value)
		case "reqINT":
			model.Required.Int = uint16(intField.Value)
		case "reqLUK":
			model.Required.Luk = uint16(intField.Value)
		case "incSTR":
			model.Ability.Str = uint16(intField.Value)
		case "incDEX":
			model.Ability.Dex = uint16(intField.Value)
		case "incINT":
			model.Ability.Int = uint16(intField.Value)
		case "incLUK":
			model.Ability.Luk = uint16(intField.Value)
		case "incPAD":
			model.Ability.PAD = uint16(intField.Value)
		case "incMAD":
			model.Ability.MAD = uint16(intField.Value)
		case "incPDD":
			model.Ability.PDD = uint16(intField.Value)
		case "incMDD":
			model.Ability.MDD = uint16(intField.Value)
		case "incPVPDamage":
			model.Ability.PVPDamage = intField.Value
		case "incSpeed":
			model.Ability.Speed = uint16(intField.Value)
		case "incJump":
			model.Ability.Jump = uint16(intField.Value)
		case "incACC":
			model.Ability.ACC = uint16(intField.Value)
		case "incEVA":
			model.Ability.EVA = intField.Value
		case "incMHP":
			model.Ability.MaxHP = uint16(intField.Value)
		case "incMMP":
			model.Ability.MaxMP = uint16(intField.Value)
		case "tuc":
			model.TUC = uint8(intField.Value)
		case "price":
			model.Price = intField.Value
		case "attackSpeed":
			model.AttackSpeed = intField.Value
		case "cash":
			model.Cash = intField.Value == 1
		case "slotMax":
			model.SlotMax = uint16(intField.Value)
		case "quest":
			model.Quest = intField.Value == 1
		case "equipTradeBlock":
			model.EquipTradeBlock = intField.Value == 1
		case "tradeAvailable":
			model.TradeAvailable = intField.Value
		case "hide":
			model.Hide = intField.Value == 1
		case "royalSpecial":
			model.RoyalSpecial = intField.Value == 1
		case "masterSpecial":
			model.MasterSpecial = intField.Value == 1
		case "walk":
		case "stand":
		case "charmEXP":
		case "bossReward":
		case "exItem":
		case "imdR":
		case "bdR":
		case "willEXP":
		case "charismaEXP":
		case "kaiserOffsetX":
		case "kaiserOffsetY":
		case "knockback":
		case "noMoveToLocker":
		case "weekly":
		case "cashTradeBlock":
		case "abilityTimeLimited":
		case "cashForceCharmExp":
		case "effect":
		case "reqPOP":
		case "invisibleFace":
		case "equippedEmotion":
		case "equippedSound":
		case "baseLevel":
		case "medalTag":
		case "tradBlock":
		case "MaxHP":
		case "bonusExp":
		case "incCraft":
		case "specialID":
		case "speed":
		case "epicItem":
		case "notExtend":
		case "origin":
		case "incMMD":
		case "accountShareTag":
		case "pachinko":
		case "incHP":
		case "onlyCash":
		case "incMHPr":
		case "incMMPr":
		case "keywordEffect":
		case "extendFrame":
		case "vehicleDefaultFrame":
		case "isAbleToTradeOnce":
		case "noExtend":
		case "incLUk":
		case "recovery":
		case "groupEffectID":
		case "sample":
		case "nameTag":
		case "chatBalloon":
		case "pickupMeso":
		case "pickupItem":
		case "pickupOthers":
		case "sweepForDrop":
		case "longRange":
		case "consumeMP":
		case "onlyEquip":
		case "scope":
		case "bestFriendPartyBonusExp":
		case "bloodAllianceExpRate":
		case "bloodAlliancePartyExpRate":
		case "fs":
		case "tamingMob":
		case "vehicleNaviFlyingLevel":
		case "vehicleDoubleJumpLevel":
		case "vehicleGlideLevel":
		case "vehicleNewFlyingLevel":
		case "vehicleSkillIsTown":
		case "passengerNum":
		case "removeBody":
		case "ActionEffect":
		case "incSwim":
		case "incFatigue":
		case "hpRecovery":
		case "mpRecovery":
		case "partsQuestID":
		case "partsCount":
		}
	}

	for _, v := range node.Children {
		switch v.Name {
		case "icon":
		case "iconRaw":
		case "islot":
		case "vslot":
		case "afterImage":
		case "sfx":
		case "tradeBlock":
		case "TradeBlock":
		case "only":
		case "timeLimited":
		case "notSale":
		case "expireOnLogout":
		case "level":
		case "setItemID":
			break

		default:
			mutex.Lock()
			if _, ok := visit[v.Name]; ok {
				mutex.Unlock()
				continue
			}
			visit[v.Name] = true
			mutex.Unlock()
			log.Printf("%s is not declared in %s\n", v.Name, filepath.Base(path))
		}
	}

	return &model, nil
}

// loadGeneralItems loads general item specifications from XML file.
func loadGeneralItems(path string) (*[]*GeneralItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*GeneralItem{}
	for _, v := range root.Children {
		model := GeneralItem{
			ItemCore: &ItemCore{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		model.ID = uint32(id)
		info := v.find("info")
		if info != nil {
			for _, intField := range info.Ints {
				switch intField.Name {
				case "price":
					model.Price = intField.Value
				case "slotMax":
					model.SlotMax = uint16(intField.Value)
				case "tradeBlock":
				case "only":
				case "timeLimited":
				case "cash":
				case "notSale":
				case "incPAD":
				case "incMAD":
				case "incACC":
				case "incEVA":
				case "incSpeed":
				case "incJump":
				case "incMaxHP":
				case "incMaxMP":
				case "incSTR":
				case "incINT":
				case "incLUK":
				case "incDEX":
				case "incReqLevel":
				case "randOption":
				case "randStat":
				case "quest":
				case "exp":
				case "grade":
				case "questId":
				case "lv":
				case "lvMin":
				case "lvMax":
				case "pquest":
				case "bigSize":
				case "pickUpBlock":
				case "showMessage":
				case "mcType":
				case "autoPrice":
				case "noDrop":
				case "notExtend":
				case "expireOnLogout":
				}
			}

			for _, iv := range info.Children {
				switch iv.Name {
				case "icon":
				case "iconRaw":
				case "iconShop":
				case "iconReward":
				case "name":
				case "uiData":
				case "message":
				case "consumeItem":
					break

				default:
					mutex.Lock()
					if _, ok := visit[iv.Name]; ok {
						mutex.Unlock()
						continue
					}
					visit[iv.Name] = true
					mutex.Unlock()
					log.Printf("%s is not declared in %s:info\n", iv.Name, filepath.Base(path))
				}
			}
		}

		specs = append(specs, &model)
	}

	return &specs, nil
}

// loadInstallations loads installation item specifications from XML file.
func loadInstallations(path string) (*[]*Installation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*Installation{}
	for _, v := range root.Children {
		model := Installation{
			ItemCore: &ItemCore{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		model.ID = uint32(id)
		info := v.find("info")
		if info != nil {
			for _, intField := range info.Ints {
				switch intField.Name {
				case "slotMax":
					model.SlotMax = uint16(intField.Value)
				case "price":
				case "tradeBlock":
				case "notSale":
				case "only":
				case "lv":
				case "timeLimited":
				case "expireOnLogout":
				case "recoveryHP":
				case "reqLevel":
				case "recoveryMP":
				case "tamingMob":
				case "reqGuildLevel":
				case "guild":
				case "bodyRelMove":
				case "accountSharable":
				case "accountShareTag":
				case "sitAction":
				case "sitEmotion":
				case "distanceX":
				case "distanceY":
				case "maxDiff":
				case "direction":
				}
			}

			for _, iv := range info.Children {
				switch iv.Name {
				case "icon":
				case "iconRaw":
				case "iconReward":
					break

				default:
					mutex.Lock()
					if _, ok := visit[iv.Name]; ok {
						mutex.Unlock()
						continue
					}
					visit[iv.Name] = true
					mutex.Unlock()
					log.Printf("%s is not declared in %s:info\n", iv.Name, filepath.Base(path))
				}
			}
		}

		specs = append(specs, &model)
	}

	return &specs, nil
}

// loadMaps loads map specifications from XML file.
func loadMaps(path string, mapId uint32) (*Map, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	model := Map{
		ID:        mapId,
		Portals:   map[uint8]Portal{},
		NpcSpawns: map[uint32]NpcSpawn{},
		MobSpawns: map[uint32]MobSpawn{},
	}

	info := root.find("info")
	if info != nil {
		for _, v := range info.Children {
			switch v.Name {
			case "mapName":
				model.Name = v.Value
			case "version":
				model.Version, _ = strconv.Atoi(v.Value)
			case "cloud":
				model.Cloud, _ = strconv.Atoi(v.Value)
			case "returnMap":
				model.ReturnMapId, _ = strconv.Atoi(v.Value)
			case "forcedReturn":
				model.ForcedReturn, _ = strconv.Atoi(v.Value)
			case "fieldLimit":
				model.FieldLimit, _ = strconv.Atoi(v.Value)
			case "VRTop":
				model.VRTop, _ = strconv.Atoi(v.Value)
			case "VRLeft":
				model.VRLeft, _ = strconv.Atoi(v.Value)
			case "VRBottom":
				model.VRBottom, _ = strconv.Atoi(v.Value)
			case "VRRight":
				model.VRRight, _ = strconv.Atoi(v.Value)
			case "hideMinimap":
				model.HideMinimap = v.Value == "1"
			case "town":
				model.IsTown = v.Value == "1"
			case "mobRate":
				f, err := strconv.ParseFloat(v.Value, 32)
				if err == nil {
					model.MobRate = float32(f)
				}
			case "bgm":
				model.BGM = v.Value
			case "mapMark":
				model.MapMark = v.Value
			case "mapDesc":
				model.MapDesc = v.Value
			case "miniMapOnOff":
				model.MiniMapOnOff = v.Value == "1"
			default:
				break
			}
		}
	}

	portals := root.find("portal")
	if portals != nil {
		for _, v := range portals.Children {
			var portal Portal

			// Parse string fields
			for _, strField := range v.Strings {
				switch strField.Name {
				case "pn":
					portal.Name = strField.Value
				case "tn":
					portal.Target = strField.Value
				case "script":
					portal.ScriptName = strField.Value
				}
			}

			// Parse int fields
			for _, intField := range v.Ints {
				switch intField.Name {
				case "pt":
					portal.Type = uint8(intField.Value)
				case "tm":
					portal.TargetMapId = int32(intField.Value)
				case "x":
					portal.Position.X = int16(intField.Value)
				case "y":
					portal.Position.Y = int16(intField.Value)
				}
			}

			// Parse children fields (for compatibility)
			for _, field := range v.Children {
				switch field.Name {
				case "pn":
					portal.Name = field.Value
				case "pt":
					if val, err := strconv.Atoi(field.Value); err == nil {
						portal.Type = uint8(val)
					}
				case "tm":
					if val, err := strconv.Atoi(field.Value); err == nil {
						portal.TargetMapId = int32(val)
					}
				case "tn":
					portal.Target = field.Value
				case "x":
					if val, err := strconv.Atoi(field.Value); err == nil {
						portal.Position.X = int16(val)
					}
				case "y":
					if val, err := strconv.Atoi(field.Value); err == nil {
						portal.Position.Y = int16(val)
					}
				case "script":
					portal.ScriptName = field.Value
				}
			}

			id, _ := strconv.Atoi(v.Name)
			portal.ID = uint8(id)
			model.Portals[portal.ID] = portal
		}
	}

	bound := types.Rect[int16]{}
	footholds := root.find("foothold")
	buffer := []Foothold{}
	if footholds != nil {
		for _, v1 := range footholds.Children {
			for _, v2 := range v1.Children {
				for _, v3 := range v2.Children {
					// Skip non-numeric nodes (e.g., "AreaCode")
					id, err := strconv.Atoi(v3.Name)
					if err != nil {
						log.Printf("Skipping non-numeric foothold node '%s' in map %d (file: %s)", v3.Name, mapId, filepath.Base(path))
						continue
					}

					// Helper function to find int field value
					findIntValue := func(fieldName string) (int, bool) {
						for _, intField := range v3.Ints {
							if intField.Name == fieldName {
								return intField.Value, true
							}
						}
						// Also check in child nodes (for compatibility)
						if childNode := v3.find(fieldName); childNode != nil {
							if val, err := strconv.Atoi(childNode.Value); err == nil {
								return val, true
							}
						}
						return 0, false
					}

					// Find and validate x1
					x1, found := findIntValue("x1")
					if !found {
						log.Printf("Missing 'x1' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					// Find and validate x2
					x2, found := findIntValue("x2")
					if !found {
						log.Printf("Missing 'x2' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					// Find and validate y1
					y1, found := findIntValue("y1")
					if !found {
						log.Printf("Missing 'y1' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					// Find and validate y2
					y2, found := findIntValue("y2")
					if !found {
						log.Printf("Missing 'y2' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					// Find and validate prev
					prev, found := findIntValue("prev")
					if !found {
						log.Printf("Missing 'prev' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					// Find and validate next
					next, found := findIntValue("next")
					if !found {
						log.Printf("Missing 'next' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					foothold := Foothold{
						ID:   int16(id),
						X1:   int16(x1),
						Y1:   int16(y1),
						X2:   int16(x2),
						Y2:   int16(y2),
						Prev: int16(prev),
						Next: int16(next),
					}
					buffer = append(buffer, foothold)

					bound.Left = min(bound.Left, foothold.X1)
					bound.Right = max(bound.Right, foothold.X2)
					bound.Top = min(bound.Top, foothold.Y1)
					bound.Bottom = max(bound.Bottom, foothold.Y2)
				}
			}
		}

		model.Footholds = types.NewQuadTreeNode[int16, Foothold](bound, 0)
		for _, foothold := range buffer {
			model.Footholds.Insert(foothold)
		}
	}

	lives := root.find("life")
	if lives != nil {
		var spawnId int
		for _, life := range lives.Children {
			spawnId, _ = strconv.Atoi(life.Name)
			var spawn Spawn
			baseSpawn := BaseSpawn{
				FacingDirection: FACING_DIRECTION_LEFT,
			}

			// Find type field - it's a <string> tag, so check Strings first
			lifeType := ""
			for _, strField := range life.Strings {
				if strField.Name == "type" {
					lifeType = strField.Value
					break
				}
			}

			// If not found in Strings, check Children (for backward compatibility)
			if lifeType == "" {
				for _, prop := range life.Children {
					if prop.Name == "type" {
						lifeType = prop.Value
						break
					}
				}
			}

			// Determine spawn type
			if lifeType == "n" {
				spawn = NpcSpawn{
					BaseSpawn: &baseSpawn,
				}
			} else if lifeType == "m" {
				spawn = MobSpawn{
					BaseSpawn: &baseSpawn,
				}
			} else {
				err := fmt.Errorf("invalid life type '%s' for life entry '%s' in map %d (file: %s)", lifeType, life.Name, mapId, filepath.Base(path))
				log.Printf("Error: %v", err)
				return nil, err
			}

			// Process string fields (id, limitedname)
			for _, strField := range life.Strings {
				switch strField.Name {
				case "id":
					value, err := strconv.Atoi(strField.Value)
					if err != nil {
						return nil, err
					}
					baseSpawn.ID = uint32(value)
				case "limitedname":
					baseSpawn.LimitedName = strField.Value
				}
			}

			// Process int fields (x, y, mobTime, f, fh, cy, rx0, rx1, hide, useDay, useNight, info, nofoothold)
			for _, intField := range life.Ints {
				switch intField.Name {
				case "x":
					baseSpawn.Position.X = int16(intField.Value)
				case "y":
					baseSpawn.Position.Y = int16(intField.Value)
				case "mobTime":
					baseSpawn.MobTime = time.Duration(intField.Value) * time.Second
				case "f":
					if intField.Value == 0 {
						baseSpawn.FacingDirection = FACING_DIRECTION_LEFT
					} else {
						baseSpawn.FacingDirection = FACING_DIRECTION_RIGHT
					}
				case "fh":
					baseSpawn.Foothold = int16(intField.Value)
				case "cy":
					baseSpawn.CollisionY = int16(intField.Value)
				case "rx0":
					baseSpawn.RenderX0 = int16(intField.Value)
				case "rx1":
					baseSpawn.RenderX1 = int16(intField.Value)
				case "hide":
					baseSpawn.Hide = intField.Value != 0
				case "useDay":
					baseSpawn.UseDay = intField.Value != 0
				case "useNight":
					baseSpawn.UseNight = intField.Value != 0
				case "info":
					baseSpawn.Info = uint8(intField.Value)
				case "nofoothold":
					baseSpawn.NoFoothold = intField.Value != 0
				}
			}

			// Process children (for backward compatibility with any remaining fields)
			for _, prop := range life.Children {
				// Skip if already processed as string or int
				alreadyProcessed := false
				for _, strField := range life.Strings {
					if strField.Name == prop.Name {
						alreadyProcessed = true
						break
					}
				}
				if !alreadyProcessed {
					for _, intField := range life.Ints {
						if intField.Name == prop.Name {
							alreadyProcessed = true
							break
						}
					}
				}
				if alreadyProcessed {
					continue
				}

				// Handle any remaining fields that might be in Children
				switch prop.Name {
				case "type":
					// Already processed above, skip
					continue
				default:
					mutex.Lock()
					if _, ok := visit[prop.Name]; ok {
						mutex.Unlock()
						continue
					}
					visit[prop.Name] = true
					mutex.Unlock()
					log.Printf("%s is not declared in %s:info\n", prop.Name, filepath.Base(path))
				}
			}

			switch v := spawn.(type) {
			case NpcSpawn:
				model.NpcSpawns[uint32(spawnId)] = v

			case MobSpawn:
				model.MobSpawns[uint32(spawnId)] = v

			default:
				err := fmt.Errorf("uninitialized spawn type for life entry '%s' in map %d (file: %s)", life.Name, mapId, filepath.Base(path))
				log.Printf("Error: %v", err)
				return nil, err
			}
			spawn = nil
		}
	}

	return &model, nil
}

// loadPets loads pet specifications from XML file.
func loadPets(path string) (*Pet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(strings.TrimSuffix(root.Name, ".img"))
	if err != nil {
		return nil, err
	}
	model := &Pet{
		ItemCore: &ItemCore{
			ID: uint32(id),
		},
	}
	info := root.find("info")
	for _, iv := range info.Children {
		switch iv.Name {
		case "mob":
		case "icon":
		case "iconRaw":
		case "iconD":
		case "iconRawD":
		case "hungry":
		case "cash":
		case "life":
		case "limitedLife":
		case "noRevive":
		case "noMoveToLocker":
		case "pickupItem":
		case "consumeHP":
		case "consumeMP":
		case "sweepForDrop":
		case "nameTag":
		case "chatBalloon":
		case "pickupAll":
		case "longRange":
		case "multiPet":
		case "autoBuff":
		case "setItemID":
			break

		default:
			mutex.Lock()
			if _, ok := visit[iv.Name]; ok {
				mutex.Unlock()
				continue
			}
			visit[iv.Name] = true
			mutex.Unlock()
			log.Printf("%s is not declared in %s:info\n", iv.Name, filepath.Base(path))
		}
	}
	return model, nil
}

// loadSpecialItems loads special item specifications from XML file.
func loadSpecialItems(path string) (*[]*SpecialItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*SpecialItem{}
	for _, v := range root.Children {
		model := SpecialItem{
			ItemCore: &ItemCore{
				SlotMax: 1,
			},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		model.ID = uint32(id)
		for _, iv := range v.Children {
			switch iv.Name {
			case "icon":
			case "name":
			case "delta":
			case "iconRaw":
			case "desc":
				break

			default:
				mutex.Lock()
				if _, ok := visit[iv.Name]; ok {
					mutex.Unlock()
					continue
				}
				visit[iv.Name] = true
				mutex.Unlock()
				log.Printf("%s is not declared in %s:info\n", iv.Name, filepath.Base(path))
			}
		}

		specs = append(specs, &model)
	}

	return &specs, nil
}

// loadStringNodeRecursive recursively loads string specifications from XML nodes.
func loadStringNodeRecursive(root *node) map[uint32]map[string]string {

	result := map[uint32]map[string]string{}
	for _, child := range root.Children {

		id, err := strconv.Atoi(child.Name)
		if err != nil {
			// If not an integer, recursively process child nodes
			childResult := loadStringNodeRecursive(&child)
			for k, v := range childResult {
				result[k] = v
			}
			continue
		}

		// Create a map for this ID's string data
		stringData := map[string]string{}
		for _, v := range child.Children {
			stringData[v.Name] = v.Value
		}

		result[uint32(id)] = stringData
	}

	return result
}

// loadStringResources loads string resources from XML file.
func loadStringResources(path string) (*map[uint32]map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	result := loadStringNodeRecursive(&root)
	return &result, nil
}

// loadMob loads monster specifications from XML file.
func loadMob(path string) (*Mob, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(strings.TrimSuffix(root.Name, ".img"))
	if err != nil {
		return nil, err
	}
	model := &Mob{ID: uint32(id)}

	info := root.find("info")
	if info == nil {
		return nil, fmt.Errorf("'info' node not found: %s", path)
	}

	for _, iv := range info.Children {
		switch iv.Name {
		case "elemAttr":
		case "PDRate":
		case "MDRate":
		case "category":
		case "firstAttack":
		case "link":
		case "skill":
		case "attack":
		case "fixedDamage":
		case "flySpeed":
		case "mpRecovery":
		case "boss":
		case "hpRecovery":
		case "removeAfter":
		case "revive":
		case "hpTagColor":
		case "hpTagBgcolor":
		case "HPgaugeHide":
		case "rareItemDropLevel":
		case "noFlip":
		case "mbookID":
		case "finalmaxHP":
		case "changeableMob":
		case "changeableMob_Type":
		case "publicReward":
		case "explosiveReward":
		case "wp":
		case "ban":
		case "chaseSpeed":
		case "default":
		case "noregen":
		case "isRemoteRange":
		case "showNotRemoteDam":
		case "ignoreMovable":
		case "ignoreMoveImpact":
		case "selfDestruction":
		case "buff":
		case "speak":
		case "useReaction":
		case "ignoreFieldOut":
		case "defaultHP":
		case "defaultMP":
		case "partyBonusMob":
		case "firstAttackRange":
		case "hideMove":
		case "mobJobCategory":
		case "invincible":
		case "hideHP":
		case "hideName":
		case "noDebuff":
		case "charismaEXP":
		case "willEXP":
		case "fixedBodyAttackDamageR":
		case "bodyDisease":
		case "bodyDiseaseLevel":
		case "summonEffect":
		case "disable":
		case "notAttack":
		case "underObject":
		case "damagedBySelectedSkill":
		case "atom":
		case "damagedByMob":
		case "dropItemPeriod":
		case "getCP":
		case "damagedBySelectedMob":
		case "doNotRemove":
		case "loseItem":
		case "thumbnail":
		case "Speed":
		case "fixDamage":
		case "nonLevelCheckACC":
		case "nonLevelCheckEVA":
		case "patrol":
		case "onlyNormalAttack":
		case "effectiveSkill":
		case "mobZone":
		case "ignoreSkill":
		case "individualReward":
		case "FlySpeed":
		case "firstattack":
		case "removeQuest":
		case "bodyattack":
		case "damageModification":

		case "bodyAttack":
			v, _ := strconv.Atoi(iv.Value)
			model.BodyAttack = v
		case "level":
			v, _ := strconv.Atoi(iv.Value)
			model.Level = uint8(v)
		case "maxHP":
			v, _ := strconv.Atoi(iv.Value)
			model.MaxHP = v
		case "maxMP":
			v, _ := strconv.Atoi(iv.Value)
			model.MaxMP = v
		case "speed":
			v, _ := strconv.Atoi(iv.Value)
			model.Speed = int16(v)
		case "PADamage":
			v, _ := strconv.Atoi(iv.Value)
			model.PADamage = v
		case "PDDamage":
			v, _ := strconv.Atoi(iv.Value)
			model.PDDamage = v
		case "MADamage":
			v, _ := strconv.Atoi(iv.Value)
			model.MADamage = v
		case "MDDamage":
			v, _ := strconv.Atoi(iv.Value)
			model.MDDamage = v
		case "acc":
			v, _ := strconv.Atoi(iv.Value)
			model.ACC = v
		case "eva":
			v, _ := strconv.Atoi(iv.Value)
			model.EVA = v
		case "exp":
			v, _ := strconv.Atoi(iv.Value)
			model.EXP = uint32(v)
		case "undead":
			v, _ := strconv.Atoi(iv.Value)
			model.Undead = (v == 1)
		case "pushed":
			v, _ := strconv.Atoi(iv.Value)
			model.Pushed = (v == 1)
		case "fs":
			f, _ := strconv.ParseFloat(iv.Value, 32)
			model.FS = float32(f)
		case "summonType":
			v, _ := strconv.Atoi(iv.Value)
			model.SummonType = uint8(v)
		case "mobType":
			v, _ := strconv.Atoi(iv.Value)
			model.MobType = uint8(v)
		default:
			mutex.Lock()
			if _, seen := visit[iv.Name]; !seen {
				visit[iv.Name] = true
				log.Printf("%s is not declared in %s:info\n", iv.Name, filepath.Base(path))
			}
			mutex.Unlock()
		}
	}

	return model, nil
}

// loadDrops loads monster drop tables from XML file.
func loadDrops(path string) (*map[uint32][]Drop, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := map[uint32][]Drop{}
	for _, v := range root.Children {
		if !strings.HasPrefix(v.Name, "m") {
			continue
		}
		id, err := strconv.Atoi(strings.TrimPrefix(v.Name, "m"))
		if err != nil {
			return nil, err
		}
		model := Drop{
			Mob: uint32(id),
		}

		for _, iv1 := range v.Children {
			for _, iv2 := range iv1.Children {
				switch iv2.Name {
				case "item":
					value, err := strconv.Atoi(iv2.Value)
					if err != nil {
						return nil, err
					}
					model.Item = uint32(value)
				case "money":
					value, err := strconv.Atoi(iv2.Value)
					if err != nil {
						return nil, err
					}
					model.Money = uint32(value)
				case "prob":
					value, err := strconv.ParseFloat(strings.TrimPrefix(iv2.Value, "[R8]"), 32)
					if err != nil {
						return nil, err
					}
					model.Prob = float32(value)
				case "min":
					value, err := strconv.Atoi(iv2.Value)
					if err != nil {
						return nil, err
					}
					model.Min = uint16(value)
				case "max":
					value, err := strconv.Atoi(iv2.Value)
					if err != nil {
						return nil, err
					}
					model.Max = uint16(value)
				default:
					mutex.Lock()
					if _, seen := visit[iv2.Name]; !seen {
						visit[iv2.Name] = true
						log.Printf("%s is not declared in %s:info\n", iv2.Name, filepath.Base(path))
					}
					mutex.Unlock()
				}
			}
			specs[model.Mob] = append(specs[model.Mob], model)
		}
	}

	return &specs, nil
}
