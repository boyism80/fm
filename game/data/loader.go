package data

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/boyism80/fm/common/types"
)

var mutex sync.Mutex = sync.Mutex{}
var visit map[string]bool = map[string]bool{}

func loadCashItems(path string) (*[]*CashItemSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*CashItemSpec{}
	for _, v := range root.Children {
		spec := CashItemSpec{
			ItemCoreSpec: &ItemCoreSpec{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		spec.ID = uint32(id)
		info := v.find("info")
		for _, iv := range info.Children {
			switch iv.Name {
			case "icon":
			case "iconRaw":
			case "cash":
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

			case "slotMax":
				slotMax, err := strconv.Atoi(iv.Value)
				if err != nil {
					return nil, err
				}
				spec.SlotMax = uint16(slotMax)

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

		specs = append(specs, &spec)
	}

	return &specs, nil
}

func loadConsumes(path string) (*[]*ConsumeSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*ConsumeSpec{}
	for _, v := range root.Children {
		spec := ConsumeSpec{
			ItemCoreSpec: &ItemCoreSpec{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		spec.ID = uint32(id)
		nodeInfo := v.find("info")
		for _, iv := range nodeInfo.Children {
			switch iv.Name {
			case "icon":
			case "iconRaw":
			case "pquest":
			case "tragetBlock":
			case "timeLimited":
			case "only":
			case "tradeBlock":
			case "notSale":
			case "quest":
			case "incPAD":
			case "success":
			case "masterLevel":
			case "skill":
			case "unitPrice":
			case "reqLevel":
			case "tradeAvailable":
			case "noCancelMouse":
			case "mob":
			case "create":
			case "left":
			case "right":
			case "top":
			case "bottom":
			case "mobHP":
			case "bridleMsgType":
			case "bridleProp":
			case "bridlePropChg":
			case "useDelay":
			case "delayMsg":
			case "reqSkillLevel":
			case "type":
			case "incPDD":
			case "incMDD":
			case "incACC":
			case "incMHP":
			case "cursed":
			case "incINT":
			case "incMAD":
			case "incDEX":
			case "incLUK":
			case "incSTR":
			case "incSpeed":
			case "incMMP":
			case "incEVA":
			case "incJump":
			case "preventslip":
			case "warmsupport":
			case "reqRUC":
			case "recover":
			case "randstat":
			case "mcType":
			case "effect":
			case "tradBlock":
			case "bigSize":
			case "scanTradeBlock":
			case "monsterBook":
				break

			case "price":
				value, err := strconv.Atoi(iv.Value)
				if err != nil {
					return nil, err
				}
				spec.Price = value

			case "slotMax":
				value, err := strconv.Atoi(iv.Value)
				if err != nil {
					return nil, err
				}
				spec.SlotMax = uint16(value)

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

		nodeSpec := v.find("spec")
		if nodeSpec != nil {
			for _, sv := range nodeSpec.Children {
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
					log.Printf("%s is not declared in %s:spec\n", sv.Name, filepath.Base(path))
				}
			}
		}

		specs = append(specs, &spec)
	}

	return &specs, nil
}

func loadWeapons(path string) (*EquipmentSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	spec := EquipmentSpec{
		ItemCoreSpec: &ItemCoreSpec{},
	}
	id, err := strconv.Atoi(strings.TrimSuffix(root.Name, ".img"))
	if err != nil {
		return nil, err
	}
	spec.ID = uint32(id)
	node := root.find("info")
	if node == nil {
		return nil, fmt.Errorf("'info' does not exist in %s", path)
	}

	for _, v := range node.Children {
		switch v.Name {
		case "icon":
		case "iconRaw":
		case "islot":
		case "vslot":
		case "walk":
		case "stand":
		case "attack":
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
			break

		case "reqJob":
			class, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Required.Class = class
		case "reqLevel":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Required.Class = value
		case "reqSTR":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Required.Str = uint16(value)
		case "reqDEX":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Required.Dex = uint16(value)
		case "reqINT":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Required.Int = uint16(value)
		case "reqLUK":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Required.Luk = uint16(value)
		case "incSTR":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.Str = uint16(value)
		case "incDEX":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.Dex = uint16(value)
		case "incINT":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.Int = uint16(value)
		case "incLUK":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.Luk = uint16(value)
		case "incPAD":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.PAD = uint16(value)
		case "incMAD":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.MAD = uint16(value)
		case "incPDD":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.PDD = uint16(value)
		case "incMDD":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.MDD = uint16(value)
		case "incPVPDamage":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.PVPDamage = value
		case "incSpeed":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.Speed = uint16(value)
		case "incJump":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.Jump = uint16(value)

		case "incACC":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.ACC = uint16(value)

		case "incEVA":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.EVA = value
		case "incMHP":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.MaxHP = uint16(value)
		case "incMMP":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Ability.MaxMP = uint16(value)
		case "tuc":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.TUC = uint8(value)
		case "price":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Price = value
		case "attackSpeed":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.AttackSpeed = value
		case "cash":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Cash = value == 1

		case "slotMax":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.SlotMax = uint16(value)

		case "quest":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Quest = value == 1

		case "equipTradeBlock":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.EquipTradeBlock = value == 1
		case "tradeAvailable":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.TradeAvailable = value
		case "hide":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.Hide = value == 1
		case "royalSpecial":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.RoyalSpecial = value == 1
		case "masterSpecial":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			spec.MasterSpecial = value == 1
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

	return &spec, nil
}

func loadGeneralItems(path string) (*[]*GeneralItemSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*GeneralItemSpec{}
	for _, v := range root.Children {
		spec := GeneralItemSpec{
			ItemCoreSpec: &ItemCoreSpec{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		spec.ID = uint32(id)
		info := v.find("info")
		for _, iv := range info.Children {
			switch iv.Name {
			case "icon":
			case "price":
			case "lvMin":
			case "lvMax":
			case "iconRaw":
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
			case "name":
			case "uiData":
			case "message":
			case "consumeItem":
			case "iconShop":
			case "lv":
			case "iconReward":
			case "pquest":
			case "bigSize":
			case "pickUpBlock":
			case "showMessage":
			case "mcType":
			case "autoPrice":
			case "noDrop":
			case "notExtend":
			case "expireOnLogout":
				break

			case "slotMax":
				slotMax, err := strconv.Atoi(iv.Value)
				if err != nil {
					return nil, err
				}
				spec.SlotMax = uint16(slotMax)

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

		specs = append(specs, &spec)
	}

	return &specs, nil
}

func loadInstallations(path string) (*[]*InstallationSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*InstallationSpec{}
	for _, v := range root.Children {
		spec := InstallationSpec{
			ItemCoreSpec: &ItemCoreSpec{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		spec.ID = uint32(id)
		info := v.find("info")
		for _, iv := range info.Children {
			switch iv.Name {
			case "price":
			case "icon":
			case "iconRaw":
			case "tradeBlock":
			case "notSale":
			case "only":
			case "lv":
			case "iconReward":
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
				break

			case "slotMax":
				slotMax, err := strconv.Atoi(iv.Value)
				if err != nil {
					return nil, err
				}
				spec.SlotMax = uint16(slotMax)

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

		specs = append(specs, &spec)
	}

	return &specs, nil
}

func loadMaps(path string, mapId uint32) (*MapSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	spec := MapSpec{
		ID:      mapId,
		Portals: map[uint8]Portal{},
		NPCs:    map[uint32]NPCSpec{},
		Mobs:    map[uint32]MobSpec{},
	}

	info := root.find("info")
	if info != nil {
		for _, v := range info.Children {
			switch v.Name {
			case "mapName":
				spec.Name = v.Value
			case "version":
				spec.Version, _ = strconv.Atoi(v.Value)
			case "cloud":
				spec.Cloud, _ = strconv.Atoi(v.Value)
			case "returnMap":
				spec.ReturnMapId, _ = strconv.Atoi(v.Value)
			case "forcedReturn":
				spec.ForcedReturn, _ = strconv.Atoi(v.Value)
			case "fieldLimit":
				spec.FieldLimit, _ = strconv.Atoi(v.Value)
			case "VRTop":
				spec.VRTop, _ = strconv.Atoi(v.Value)
			case "VRLeft":
				spec.VRLeft, _ = strconv.Atoi(v.Value)
			case "VRBottom":
				spec.VRBottom, _ = strconv.Atoi(v.Value)
			case "VRRight":
				spec.VRRight, _ = strconv.Atoi(v.Value)
			case "hideMinimap":
				spec.HideMinimap = v.Value == "1"
			case "town":
				spec.IsTown = v.Value == "1"
			case "mobRate":
				f, err := strconv.ParseFloat(v.Value, 32)
				if err == nil {
					spec.MobRate = float32(f)
				}
			case "bgm":
				spec.BGM = v.Value
			case "mapMark":
				spec.MapMark = v.Value
			case "mapDesc":
				spec.MapDesc = v.Value
			case "miniMapOnOff":
				spec.MiniMapOnOff = v.Value == "1"
			default:
				break
			}
		}
	}

	portals := root.find("portal")
	if portals != nil {
		for _, v := range portals.Children {
			var portal Portal
			for _, field := range v.Children {
				switch field.Name {
				case "pn":
					portal.Name = field.Value
				case "pt":
					v, _ := strconv.Atoi(field.Value)
					portal.Type = uint8(v)
				case "tm":
					v, _ := strconv.Atoi(field.Value)
					portal.TargetMapId = int32(v)
				case "tn":
					portal.Target = field.Value
				case "x":
					v, _ := strconv.Atoi(field.Value)
					portal.Position.X = int16(v)
				case "y":
					v, _ := strconv.Atoi(field.Value)
					portal.Position.Y = int16(v)
				case "script":
					if field.Value != "" {
						portal.ScriptName = field.Value
					}
				}
			}
			id, _ := strconv.Atoi(v.Name)
			portal.ID = uint8(id)
			spec.Portals[portal.ID] = portal
		}
	}

	bound := types.Rect[int16]{}
	footholds := root.find("foothold")
	buffer := []Foothold{}
	if footholds != nil {
		for _, v1 := range footholds.Children {
			for _, v2 := range v1.Children {
				for _, v3 := range v2.Children {
					id, err := strconv.Atoi(v3.Name)
					if err != nil {
						continue
					}

					x1, err := strconv.Atoi(v3.find("x1").Value)
					if err != nil {
						continue
					}

					x2, err := strconv.Atoi(v3.find("x2").Value)
					if err != nil {
						continue
					}

					y1, err := strconv.Atoi(v3.find("y1").Value)
					if err != nil {
						continue
					}

					y2, err := strconv.Atoi(v3.find("y2").Value)
					if err != nil {
						continue
					}

					prev, err := strconv.Atoi(v3.find("prev").Value)
					if err != nil {
						continue
					}

					next, err := strconv.Atoi(v3.find("next").Value)
					if err != nil {
						continue
					}

					foothold := Foothold{
						ID:   id,
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

		spec.Footholds = types.NewQuadTreeNode[int16, Foothold](bound, 0)
		for _, foothold := range buffer {
			spec.Footholds.Insert(foothold)
		}
	}

	lives := root.find("life")
	if lives != nil {
		for _, life := range lives.Children {
			var createdSpec Life
			baseSpec := LifeSpec{}
			for _, prop := range life.Children {
				switch prop.Name {
				case "type":
					if prop.Value == "n" {
						createdSpec = NPCSpec{
							LifeSpec: &baseSpec,
						}
					} else if prop.Value == "m" {
						createdSpec = MobSpec{
							LifeSpec: &baseSpec,
						}
					} else {
						return nil, errors.New("invalid life type")
					}

				case "id":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.ID = uint32(value)
				case "x":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.Position.X = int16(value)
				case "y":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.Position.Y = int16(value)
				case "mobTime":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.MobTime = uint64(value)
				case "f":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					if value == 0 {
						baseSpec.FacingDirection = FACING_DIRECTION_RIGHT
					} else {
						baseSpec.FacingDirection = FACING_DIRECTION_LEFT
					}
				case "fh":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.Foothold = int16(value)
				case "cy":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.CollisionY = int16(value)
				case "rx0":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.RenderX0 = int16(value)
				case "rx1":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.RenderX1 = int16(value)
				case "hide":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.Hide = value != 0
				case "useDay":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.UseDay = value != 0
				case "useNight":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.UseNight = value != 0
				case "limitedname":
					baseSpec.LimitedName = prop.Value
				case "info":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.Info = uint8(value)
				case "nofoothold":
					value, err := strconv.Atoi(prop.Value)
					if err != nil {
						return nil, err
					}
					baseSpec.NoFoothold = value != 0

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

			switch v := createdSpec.(type) {
			case NPCSpec:
				spec.NPCs[v.ID] = v

			case MobSpec:
				spec.Mobs[v.ID] = v

			default:
				return nil, errors.New("invalid life type")
			}
			createdSpec = nil
		}
	}

	return &spec, nil
}

func loadPets(path string) (*PetSpec, error) {
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
	spec := &PetSpec{
		ItemCoreSpec: &ItemCoreSpec{
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
	return spec, nil
}

func loadSpecialItems(path string) (*[]*SpecialItemSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*SpecialItemSpec{}
	for _, v := range root.Children {
		spec := SpecialItemSpec{
			ItemCoreSpec: &ItemCoreSpec{
				SlotMax: 1,
			},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		spec.ID = uint32(id)
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

		specs = append(specs, &spec)
	}

	return &specs, nil
}

func loadStringNodeRecursive(root *node) []*StringSpec {

	templates := []*StringSpec{}
	for _, child := range root.Children {

		id, err := strconv.Atoi(child.Name)
		if err != nil {
			templates = append(templates, loadStringNodeRecursive(&child)...)
		}
		spec := StringSpec{
			ID: uint32(id),
		}
		for _, v := range child.Children {
			switch v.Name {
			case "name":
				spec.Name = v.Value
			case "desc":
				spec.Desc = v.Value
			}
		}
	}

	return templates
}

func loadStringResources(path string) (*[]*StringSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	templates := loadStringNodeRecursive(&root)
	return &templates, nil
}
