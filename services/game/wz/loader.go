package wz

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

var mutex sync.Mutex = sync.Mutex{}
var visit map[string]bool = map[string]bool{}

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
			MoveTo:   -1,
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
					model.ScrollIncPAD = int16(intField.Value)
				case "reqLevel":
				case "tradeAvailable":
				case "incPDD":
					model.ScrollIncPDD = int16(intField.Value)
				case "incMDD":
					model.ScrollIncMDD = int16(intField.Value)
				case "incACC":
					model.ScrollIncACC = int16(intField.Value)
				case "incMHP":
					model.ScrollIncMaxHP = int16(intField.Value)
				case "incINT":
					model.ScrollIncInt = int16(intField.Value)
				case "incMAD":
					model.ScrollIncMAD = int16(intField.Value)
				case "incDEX":
					model.ScrollIncDex = int16(intField.Value)
				case "incLUK":
					model.ScrollIncLuk = int16(intField.Value)
				case "incSTR":
					model.ScrollIncStr = int16(intField.Value)
				case "incSpeed":
					model.ScrollIncSpeed = int16(intField.Value)
				case "incMMP":
					model.ScrollIncMaxMP = int16(intField.Value)
				case "incEVA":
					model.ScrollIncAvoid = int16(intField.Value)
				case "incJump":
					model.ScrollIncJump = int16(intField.Value)
				case "incCraft":
					model.ScrollIncHands = int16(intField.Value)
				case "tradBlock":
				case "bigSize":
				case "scanTradeBlock":
				case "mcType":
				case "cursed":
					model.ScrollCursed = int32(intField.Value)
				case "preventslip":
				case "warmsupport":
				case "reqRUC":
				case "recover":
					model.ScrollRecover = int32(intField.Value)
				case "randstat":
					model.ScrollRandStat = int32(intField.Value)
				case "unitPrice":
				case "useDelay":
				case "delayMsg":
				case "reqSkillLevel":
				case "success":
					model.ScrollSuccess = int32(intField.Value)
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

		specNode := v.find("spec")
		if specNode != nil {
			for _, intField := range specNode.Ints {
				switch intField.Name {
				case "hp":
					model.ActiveEffect.HP = intField.Value
				case "mp":
					model.ActiveEffect.MP = intField.Value
				case "hpR":
					model.ActiveEffect.HPRate = intField.Value
				case "mpR":
					model.ActiveEffect.MPRate = intField.Value
				case "time":
					model.BuffDuration = time.Duration(intField.Value) * time.Millisecond
				case "consumeOnPickup":
					model.ConsumeOnPickup = intField.Value != 0
				case "party":
					model.Party = intField.Value != 0
				case "moveTo":
					model.MoveTo = int32(intField.Value)
				case "expinc":
					model.ExpInc = int32(intField.Value)
				default:
					if debuffFlag, ok := consumeSpecKeyToCureDebuffFlag(intField.Name); ok && intField.Value > 0 {
						model.CureDebuffs = append(model.CureDebuffs, debuffFlag)
					} else if buffFlag, ok := consumeSpecKeyToBuffFlag(intField.Name); ok {
						if model.BuffValues == nil {
							model.BuffValues = make(map[constant.BuffFlag]int32)
						}
						model.BuffValues[buffFlag] = int32(intField.Value)
					}
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

		if model.SlotMax == 0 {
			model.SlotMax = 100
		}

		specs = append(specs, &model)
	}

	return &specs, nil
}

func consumeSpecKeyToCureDebuffFlag(specKey string) (constant.DebuffFlag, bool) {
	switch specKey {
	case "poison":
		return constant.DebuffFlagPoison, true
	case "seal":
		return constant.DebuffFlagSeal, true
	case "darkness":
		return constant.DebuffFlagDarkness, true
	case "weakness":
		return constant.DebuffFlagWeaken, true
	case "curse":
		return constant.DebuffFlagCurse, true
	default:
		return constant.DebuffFlag{}, false
	}
}

func consumeSpecKeyToBuffFlag(specKey string) (constant.BuffFlag, bool) {
	switch specKey {
	case "acc":
		return constant.BuffFlagAcc, true
	case "eva":
		return constant.BuffFlagAvoid, true
	case "jump":
		return constant.BuffFlagJump, true
	case "mad":
		return constant.BuffFlagMagicAtk, true
	case "mdd":
		return constant.BuffFlagMagicDef, true
	case "morph":
		return constant.BuffFlagMorph, true
	case "pad":
		return constant.BuffFlagWeaponAtk, true
	case "pdd":
		return constant.BuffFlagWeaponDef, true
	case "speed":
		return constant.BuffFlagSpeed, true
	case "itemupbyitem":
		return constant.BuffFlagDropRate, true
	case "mesoupbyitem":
		return constant.BuffFlagMesoRate, true
	default:
		return constant.BuffFlag{}, false
	}
}

func loadWeapons(path string) (Item, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	model := EquipmentCore{
		ItemCore: &ItemCore{},
	}
	id, err := strconv.Atoi(strings.TrimSuffix(root.Name, ".img"))
	if err != nil {

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
			model.EnhanceChance = uint8(intField.Value)
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

	if model.SlotMax == 0 {
		model.SlotMax = 1
	}

	eq := &model
	if constant.GetEquipmentType(model.ID) == constant.EquipmentTypeWeapon {
		return &Weapon{EquipmentCore: eq}, nil
	}
	return &Armor{EquipmentCore: eq}, nil
}

func loadMiscItems(path string) (*[]*MiscItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	specs := []*MiscItem{}
	for _, v := range root.Children {
		model := MiscItem{
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

		if model.SlotMax == 0 {
			model.SlotMax = 100
		}

		specs = append(specs, &model)
	}

	return &specs, nil
}

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

		if model.SlotMax == 0 {
			model.SlotMax = 100
		}

		specs = append(specs, &model)
	}

	return &specs, nil
}

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
		ID:            mapId,
		Portals:       map[uint8]Portal{},
		NpcSpawns:     map[uint32]NpcSpawn{},
		MobSpawns:     map[uint32]MobSpawn{},
		ReactorSpawns: map[uint32]ReactorSpawn{},
	}

	info := root.find("info")
	if info != nil {
		for _, iv := range info.Ints {
			switch iv.Name {
			case "version":
				model.Version = iv.Value
			case "cloud":
				model.Cloud = iv.Value
			case "returnMap":
				model.ReturnMapId = iv.Value
			case "forcedReturn":
				model.ForcedReturn = iv.Value
			case "fieldLimit":
				model.FieldLimit = iv.Value
			case "VRTop":
				model.VRTop = iv.Value
			case "VRLeft":
				model.VRLeft = iv.Value
			case "VRBottom":
				model.VRBottom = iv.Value
			case "VRRight":
				model.VRRight = iv.Value
			case "hideMinimap":
				model.HideMinimap = iv.Value != 0
			case "town":
				model.IsTown = iv.Value == 1
			case "everlast":
				model.Everlast = iv.Value > 0
			case "mobRate":
				model.MobRate = float32(iv.Value)
			case "recoveryRate":
				if iv.Value > 0 {
					model.RecoveryRate = float32(iv.Value)
				} else {
					model.RecoveryRate = 1.0
				}
			case "miniMapOnOff":
				model.MiniMapOnOff = iv.Value != 0
			}
		}
		for _, fv := range info.Floats {
			if f, err := strconv.ParseFloat(fv.Value, 32); err == nil {
				switch fv.Name {
				case "mobRate":
					model.MobRate = float32(f)
				case "recoveryRate":
					if f > 0 {
						model.RecoveryRate = float32(f)
					}
				}
			}
		}
		for _, sv := range info.Strings {
			switch sv.Name {
			case "bgm":
				model.BGM = sv.Value
			case "mapMark":
				model.MapMark = sv.Value
			case "mapDesc":
				model.MapDesc = sv.Value
			}
		}
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
			case "everlast":
				model.Everlast = v.Value != "0" && v.Value != ""
			case "mobRate":
				f, err := strconv.ParseFloat(v.Value, 32)
				if err == nil {
					model.MobRate = float32(f)
				}
			case "recoveryRate":
				f, err := strconv.ParseFloat(v.Value, 32)
				if err == nil && f > 0 {
					model.RecoveryRate = float32(f)
				} else {
					model.RecoveryRate = 1.0
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
		nextDoorPortalID := 0x80
		for _, v := range portals.Children {
			var portal Portal

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

			var portalID uint8
			if portal.Type == 6 {
				if nextDoorPortalID > 0xff {
					log.Printf("door portal id overflow (>=256) in map %d (file: %s)", mapId, filepath.Base(path))
					continue
				}
				portalID = uint8(nextDoorPortalID)
				nextDoorPortalID++
			} else {
				id, err := strconv.Atoi(v.Name)
				if err != nil {
					log.Printf("Skipping portal with non-numeric name '%s' in map %d (file: %s)", v.Name, mapId, filepath.Base(path))
					continue
				}
				if id < 0 || id > 255 {
					log.Printf("Portal id %d out of uint8 range in map %d (file: %s)", id, mapId, filepath.Base(path))
					continue
				}
				portalID = uint8(id)
			}
			portal.ID = portalID
			model.Portals[portal.ID] = portal
		}
	}
	model.buildDoorReturnPortal()

	if areaNode := root.find("area"); areaNode != nil {
		type areaEntry struct {
			idx  int
			rect types.Rect[int16]
		}
		entries := make([]areaEntry, 0, len(areaNode.Children))
		for _, child := range areaNode.Children {
			idx, err := strconv.Atoi(child.Name)
			if err != nil || idx < 0 {
				continue
			}
			var x1, y1, x2, y2 int
			var foundX1, foundY1, foundX2, foundY2 bool
			for _, intField := range child.Ints {
				switch intField.Name {
				case "x1":
					x1, foundX1 = intField.Value, true
				case "y1":
					y1, foundY1 = intField.Value, true
				case "x2":
					x2, foundX2 = intField.Value, true
				case "y2":
					y2, foundY2 = intField.Value, true
				}
			}
			for _, field := range child.Children {
				switch field.Name {
				case "x1":
					if val, err := strconv.Atoi(field.Value); err == nil {
						x1, foundX1 = val, true
					}
				case "y1":
					if val, err := strconv.Atoi(field.Value); err == nil {
						y1, foundY1 = val, true
					}
				case "x2":
					if val, err := strconv.Atoi(field.Value); err == nil {
						x2, foundX2 = val, true
					}
				case "y2":
					if val, err := strconv.Atoi(field.Value); err == nil {
						y2, foundY2 = val, true
					}
				}
			}
			if !foundX1 || !foundY1 || !foundX2 || !foundY2 {
				continue
			}
			entries = append(entries, areaEntry{
				idx: idx,
				rect: types.Rect[int16]{
					Left:   int16(x1),
					Top:    int16(y1),
					Right:  int16(x2),
					Bottom: int16(y2),
				},
			})
		}
		if len(entries) > 0 {
			sort.Slice(entries, func(i, j int) bool {
				return entries[i].idx < entries[j].idx
			})
			maxIdx := entries[len(entries)-1].idx
			model.Areas = make([]types.Rect[int16], maxIdx+1)
			for _, e := range entries {
				model.Areas[e.idx] = e.rect
			}
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
						log.Printf("Skipping non-numeric foothold node '%s' in map %d (file: %s)", v3.Name, mapId, filepath.Base(path))
						continue
					}

					findIntValue := func(fieldName string) (int, bool) {
						for _, intField := range v3.Ints {
							if intField.Name == fieldName {
								return intField.Value, true
							}
						}

						if childNode := v3.find(fieldName); childNode != nil {
							if val, err := strconv.Atoi(childNode.Value); err == nil {
								return val, true
							}
						}
						return 0, false
					}

					x1, found := findIntValue("x1")
					if !found {
						log.Printf("Missing 'x1' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					x2, found := findIntValue("x2")
					if !found {
						log.Printf("Missing 'x2' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					y1, found := findIntValue("y1")
					if !found {
						log.Printf("Missing 'y1' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					y2, found := findIntValue("y2")
					if !found {
						log.Printf("Missing 'y2' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

					prev, found := findIntValue("prev")
					if !found {
						log.Printf("Missing 'prev' field in foothold node '%s' (ID: %d) in map %d (file: %s)", v3.Name, id, mapId, filepath.Base(path))
						continue
					}

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

			lifeType := ""
			for _, strField := range life.Strings {
				if strField.Name == "type" {
					lifeType = strField.Value
					break
				}
			}

			if lifeType == "" {
				for _, prop := range life.Children {
					if prop.Name == "type" {
						lifeType = prop.Value
						break
					}
				}
			}

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

			for _, prop := range life.Children {

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

				switch prop.Name {
				case "type":

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

	reactors := root.find("reactor")
	if reactors != nil {
		for _, reactorNode := range reactors.Children {
			spawnId, err := strconv.Atoi(reactorNode.Name)
			if err != nil {
				log.Printf("Skipping non-numeric reactor node '%s' in map %d (file: %s)", reactorNode.Name, mapId, filepath.Base(path))
				continue
			}

			spawn := ReactorSpawn{}
			for _, strField := range reactorNode.Strings {
				switch strField.Name {
				case "id":
					value, parseErr := strconv.Atoi(strField.Value)
					if parseErr != nil {
						return nil, parseErr
					}
					spawn.ReactorID = uint32(value)
				case "name":
					spawn.Name = strField.Value
				}
			}
			for _, intField := range reactorNode.Ints {
				switch intField.Name {
				case "x":
					spawn.Position.X = int16(intField.Value)
				case "y":
					spawn.Position.Y = int16(intField.Value)
				case "f":
					if intField.Value == 0 {
						spawn.FacingDirection = FACING_DIRECTION_LEFT
					} else {
						spawn.FacingDirection = FACING_DIRECTION_RIGHT
					}
				case "reactorTime":
					if intField.Value > 0 {
						spawn.RespawnDelay = time.Duration(intField.Value) * time.Second
					}
				}
			}
			for _, field := range reactorNode.Children {
				switch field.Name {
				case "id":
					if spawn.ReactorID == 0 {
						value, parseErr := strconv.Atoi(field.Value)
						if parseErr != nil {
							return nil, parseErr
						}
						spawn.ReactorID = uint32(value)
					}
				case "name":
					if spawn.Name == "" {
						spawn.Name = field.Value
					}
				case "x":
					if val, parseErr := strconv.Atoi(field.Value); parseErr == nil {
						spawn.Position.X = int16(val)
					}
				case "y":
					if val, parseErr := strconv.Atoi(field.Value); parseErr == nil {
						spawn.Position.Y = int16(val)
					}
				case "f":
					if val, parseErr := strconv.Atoi(field.Value); parseErr == nil {
						if val == 0 {
							spawn.FacingDirection = FACING_DIRECTION_LEFT
						} else {
							spawn.FacingDirection = FACING_DIRECTION_RIGHT
						}
					}
				case "reactorTime":
					if val, parseErr := strconv.Atoi(field.Value); parseErr == nil && val > 0 {
						spawn.RespawnDelay = time.Duration(val) * time.Second
					}
				}
			}

			if spawn.ReactorID == 0 {
				log.Printf("Skipping reactor spawn %d with missing id in map %d (file: %s)", spawnId, mapId, filepath.Base(path))
				continue
			}
			model.ReactorSpawns[uint32(spawnId)] = spawn
		}
	}

	return &model, nil
}

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
			ID:      uint32(id),
			SlotMax: 100,
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

func loadStringNodeRecursive(root *node) map[uint32]map[string]string {

	result := map[uint32]map[string]string{}
	for _, child := range root.Children {

		id, err := strconv.Atoi(child.Name)
		if err != nil {

			childResult := loadStringNodeRecursive(&child)
			for k, v := range childResult {
				result[k] = v
			}
			continue
		}

		stringData := map[string]string{}
		for _, v := range child.Children {
			stringData[v.Name] = v.Value
		}

		result[uint32(id)] = stringData
	}

	return result
}

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

func parseUnitPrice(unitPriceStr string) float64 {

	priceStr := strings.TrimPrefix(unitPriceStr, "[R8]")
	priceStr = strings.TrimPrefix(priceStr, "[R4]")
	priceStr = strings.Trim(priceStr, "[]")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0
	}
	return price
}

func loadNpcShops(path string) (*map[uint32]*Shop, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root shopRoot
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	result := make(map[uint32]*Shop)

	for _, shopNode := range root.Shops {
		npcID, err := strconv.ParseUint(shopNode.Name, 10, 32)
		if err != nil {
			continue
		}

		shop := &Shop{
			NpcID: uint32(npcID),
			Items: make([]ShopItem, 0, len(shopNode.Items)),
		}

		for _, itemNode := range shopNode.Items {
			item := ShopItem{}

			for _, intField := range itemNode.Ints {
				switch intField.Name {
				case "item":
					item.ItemID = uint32(intField.Value)
				case "price":
					item.Price = intField.Value
				case "period":
					item.Period = intField.Value
				case "stock":
					item.Stock = intField.Value
				}
			}

			for _, strField := range itemNode.Strings {
				if strField.Name == "unitPrice" {
					item.UnitPrice = parseUnitPrice(strField.Value)
				}
			}

			if item.ItemID > 0 {
				if item.ItemID/10000 == constant.ItemCategoryShuriken && item.ItemID != constant.ItemShurikenBase {
					continue
				}
				shop.Items = append(shop.Items, item)
			}
		}

		rechargeableItems := make([]uint32, 0, len(constant.RechargeableShurikens)+len(constant.RechargeableBullets))
		rechargeableItems = append(rechargeableItems, constant.RechargeableShurikens...)
		rechargeableItems = append(rechargeableItems, constant.RechargeableBullets...)

		existingRechargeable := make(map[uint32]bool)
		for _, existingItem := range shop.Items {
			category := existingItem.ItemID / 10000
			if category == constant.ItemCategoryShuriken || category == constant.ItemCategoryBullet {
				existingRechargeable[existingItem.ItemID] = true
			}
		}

		for _, rechargeID := range rechargeableItems {
			if !existingRechargeable[rechargeID] {
				shop.Items = append(shop.Items, ShopItem{
					ItemID:    rechargeID,
					Price:     0,
					Period:    0,
					Stock:     0,
					UnitPrice: 0,
				})
			}
		}

		if len(shop.Items) > 0 {
			result[shop.NpcID] = shop
		}
	}

	return &result, nil
}

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

	for _, intField := range info.Ints {
		switch intField.Name {
		case "bodyAttack":
			model.BodyAttack = intField.Value
		case "level":
			model.Level = uint8(intField.Value)
		case "maxHP":
			model.MaxHP = intField.Value
		case "maxMP":
			model.MaxMP = intField.Value
		case "speed":
			model.Speed = int16(intField.Value)
		case "PADamage":
			model.PADamage = intField.Value
		case "PDDamage":
			model.PDDamage = intField.Value
		case "MADamage":
			model.MADamage = intField.Value
		case "MDDamage":
			model.MDDamage = intField.Value
		case "acc":
			model.ACC = intField.Value
		case "eva":
			model.EVA = intField.Value
		case "exp":
			model.EXP = uint32(intField.Value)
		case "undead":
			model.Undead = (intField.Value == 1)
		case "pushed":
			model.Pushed = (intField.Value == 1)
		case "summonType":
			model.SummonType = uint8(intField.Value)
		case "mobType":
			model.MobType = uint8(intField.Value)
		case "boss":
			model.Boss = (intField.Value != 0)
		case "publicReward":
			model.FfaLoot = intField.Value > 0
		case "explosiveReward":
			model.ExplosiveReward = intField.Value > 0
		case "removeAfter":
			model.RemoveAfter = intField.Value
		case "hpTagColor":
			model.HpTagColor = uint8(intField.Value)
		case "hpTagBgcolor":
			model.HpTagBgColor = uint8(intField.Value)
		case "dropItemPeriod":
			model.DropItemPeriod = intField.Value
		case "damagedByMob":
			model.DamagedByMob = intField.Value > 0
		}
	}

	for _, strField := range info.Strings {
		switch strField.Name {
		case "link":
			model.Link = strField.Value
		}
	}

	for _, iv := range info.Children {
		switch iv.Name {
		case "elemAttr":
			for _, el := range iv.Children {
				key := strings.ToLower(el.Name)
				for _, inf := range el.Ints {
					if inf.Name == "value" {
						if model.ElemResist == nil {
							model.ElemResist = make(map[string]int)
						}
						model.ElemResist[key] = inf.Value
					}
				}
			}
		case "PDRate":
		case "MDRate":
		case "category":
		case "firstAttack":
		case "link":
		case "skill":
			model.Skills = parseMobInfoSkills(iv)
		case "attack":
		case "fixedDamage":
		case "flySpeed":
		case "mpRecovery":
		case "boss":
			model.Boss = true
		case "hpRecovery":
		case "removeAfter":
			model.RemoveAfter = nodeInt(&iv, "removeAfter", model.RemoveAfter)
		case "revive":
			for _, child := range iv.Children {
				for _, intf := range child.Ints {
					if intf.Value > 0 {
						model.Revives = append(model.Revives, uint32(intf.Value))
					}
				}
			}
			for _, intf := range iv.Ints {
				if intf.Value > 0 {
					model.Revives = append(model.Revives, uint32(intf.Value))
				}
			}
		case "hpTagColor":
			model.HpTagColor = uint8(nodeInt(&iv, "hpTagColor", int(model.HpTagColor)))
		case "hpTagBgcolor":
			model.HpTagBgColor = uint8(nodeInt(&iv, "hpTagBgcolor", int(model.HpTagBgColor)))
		case "HPgaugeHide":
		case "rareItemDropLevel":
		case "noFlip":
		case "mbookID":
		case "finalmaxHP":
		case "changeableMob":
		case "changeableMob_Type":
		case "publicReward":
			if v, err := strconv.Atoi(iv.Value); err == nil {
				model.FfaLoot = v > 0
			}
		case "explosiveReward":
			if v, err := strconv.Atoi(iv.Value); err == nil {
				model.ExplosiveReward = v > 0
			}
		case "wp":
		case "ban":
			model.Banish = parseMobBanish(&iv)
		case "chaseSpeed":
		case "default":
		case "noregen":
		case "isRemoteRange":
		case "showNotRemoteDam":
		case "ignoreMovable":
		case "ignoreMoveImpact":
		case "selfDestruction":
			for _, sdf := range iv.Ints {
				switch sdf.Name {
				case "removeAfter":
					model.RemoveAfter = sdf.Value
				case "action":
					model.SelfDestructionAction = int8(sdf.Value)
				}
			}
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
		case "noMobStatus":
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

		case "fs":
			f, _ := strconv.ParseFloat(iv.Value, 32)
			model.FS = float32(f)
		default:
			mutex.Lock()
			if _, seen := visit[iv.Name]; !seen {
				visit[iv.Name] = true
				log.Printf("%s is not declared in %s:info\n", iv.Name, filepath.Base(path))
			}
			mutex.Unlock()
		}
	}

	model.Attacks = parseMobAttacks(root)

	return model, nil
}

type dropEntry struct {
	XMLName xml.Name      `xml:"imgdir"`
	Name    string        `xml:"name,attr"`
	Ints    []intField    `xml:"int"`
	Strings []stringField `xml:"string"`
}

type mobDropNode struct {
	XMLName xml.Name    `xml:"imgdir"`
	Name    string      `xml:"name,attr"`
	Entries []dropEntry `xml:"imgdir"`
}

type rewardRoot struct {
	XMLName xml.Name      `xml:"imgdir"`
	Name    string        `xml:"name,attr"`
	Mobs    []mobDropNode `xml:"imgdir"`
}

type shopItemNode struct {
	XMLName xml.Name      `xml:"imgdir"`
	Name    string        `xml:"name,attr"`
	Ints    []intField    `xml:"int"`
	Strings []stringField `xml:"string"`
}

type shopNode struct {
	XMLName xml.Name       `xml:"imgdir"`
	Name    string         `xml:"name,attr"`
	Items   []shopItemNode `xml:"imgdir"`
}

type shopRoot struct {
	XMLName xml.Name   `xml:"imgdir"`
	Name    string     `xml:"name,attr"`
	Shops   []shopNode `xml:"imgdir"`
}

func parseDropEntry(entry dropEntry) Drop {
	model := Drop{}

	for _, intField := range entry.Ints {
		switch intField.Name {
		case "item":
			model.Item = uint32(intField.Value)
		case "money":
			model.Money = uint32(intField.Value)
		case "min":
			model.Min = uint16(intField.Value)
		case "max":
			model.Max = uint16(intField.Value)
		case "quest", "questid":
			model.QuestID = uint32(intField.Value)
		}
	}

	for _, stringField := range entry.Strings {
		if stringField.Name == "prob" {
			probStr := strings.TrimPrefix(stringField.Value, "[R8]")
			value, err := strconv.ParseFloat(probStr, 32)
			if err == nil {
				model.Prob = float32(value)
			}
		}
	}

	return model
}

func loadRewardDrops(path string) (map[uint32][]Drop, map[uint32][]Drop, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var root rewardRoot
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, nil, err
	}

	mobDrops := map[uint32][]Drop{}
	reactorDrops := map[uint32][]Drop{}
	for _, node := range root.Mobs {
		if strings.HasPrefix(node.Name, "m") {
			id, err := strconv.Atoi(strings.TrimPrefix(node.Name, "m"))
			if err != nil {
				continue
			}
			mobID := uint32(id)
			for _, entry := range node.Entries {
				mobDrops[mobID] = append(mobDrops[mobID], parseDropEntry(entry))
			}
		} else if strings.HasPrefix(node.Name, "r") {
			id, err := strconv.Atoi(strings.TrimPrefix(node.Name, "r"))
			if err != nil {
				continue
			}
			reactorID := uint32(id)
			for _, entry := range node.Entries {
				reactorDrops[reactorID] = append(reactorDrops[reactorID], parseDropEntry(entry))
			}
		}
	}

	return mobDrops, reactorDrops, nil
}

func loadExpTable(path string) ([]uint32, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	info := root.find("info")
	var expNode *node
	if info != nil {
		expNode = info.find("exp")
	}

	if expNode == nil {
		expNode = &root
	}

	expTable := make([]uint32, 201)

	if expNode != nil {
		for _, child := range expNode.Children {
			level, err := strconv.Atoi(child.Name)
			if err != nil {
				continue
			}
			if level < 0 || level > 200 {
				continue
			}

			var expValue int
			for _, intField := range child.Ints {
				if intField.Name == "" || intField.Name == "value" {
					expValue = intField.Value
					break
				}
			}
			if expValue == 0 {
				for _, intField := range child.Ints {
					expValue = intField.Value
					break
				}
			}
			if expValue == 0 && child.Value != "" {
				if val, err := strconv.Atoi(child.Value); err == nil {
					expValue = val
				}
			}

			if expValue > 0 {
				expTable[level] = uint32(expValue)
			}
		}
	}

	hasData := false
	for i := 1; i <= 10; i++ {
		if expTable[i] > 0 {
			hasData = true
			break
		}
	}

	if !hasData {
		return nil, fmt.Errorf("no exp table data found in %s", path)
	}

	return expTable, nil
}

func loadSkillClassFile(path string) (map[uint32]*Skill, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	skills := make(map[uint32]*Skill)

	skillNode := root.find("skill")
	if skillNode == nil {
		return skills, nil
	}

	for _, skillChild := range skillNode.Children {
		skillIDStr := skillChild.Name
		skillID, err := strconv.Atoi(skillIDStr)
		if err != nil {
			continue
		}

		skill := &Skill{
			ID:        uint32(skillID),
			LevelData: make(map[int]*SkillLevelData),
		}

		for _, strField := range skillChild.Strings {
			switch strField.Name {
			case "elemAttr":
				skill.ElemAttr = strField.Value
			}
		}

		for _, intField := range skillChild.Ints {
			switch intField.Name {
			case "masterLevel":
				skill.MasterLevel = intField.Value
			case "invisible":
				if intField.Value > 0 {
					skill.Invisible = true
				}
			case "timeLimited":
				if intField.Value > 0 {
					skill.TimeLimited = true
				}
			case "combatOrders":
				if intField.Value > 0 {
					skill.CombatOrders = true
				}
			}
		}

		info := skillChild.find("info")
		if info != nil {
			for _, intField := range info.Ints {
				switch intField.Name {
				case "masterLevel":
					if skill.MasterLevel == 0 {
						skill.MasterLevel = intField.Value
					}
				case "invisible":
					if intField.Value > 0 {
						skill.Invisible = true
					}
				case "timeLimited":
					if intField.Value > 0 {
						skill.TimeLimited = true
					}
				case "combatOrders":
					if intField.Value > 0 {
						skill.CombatOrders = true
					}
				}
			}
		}

		common := skillChild.find("common")
		if common != nil {
			for _, intField := range common.Ints {
				if intField.Name == "maxLevel" {
					skill.MaxLevel = intField.Value
					break
				}
			}
		}

		levelNode := skillChild.find("level")
		if levelNode != nil {
			if skill.MaxLevel == 0 {
				skill.MaxLevel = len(levelNode.Children)
			}

			for _, levelChild := range levelNode.Children {
				levelNum, err := strconv.Atoi(levelChild.Name)
				if err != nil {
					continue
				}

				levelData := &SkillLevelData{}

				for _, intField := range levelChild.Ints {
					switch intField.Name {
					case "mpCon":
						levelData.MPCon = intField.Value
					case "hpCon":
						levelData.HPCon = intField.Value
					case "moneyCon":
						levelData.MoneyCon = intField.Value
					case "itemCon":
						levelData.ItemCon = intField.Value
					case "itemConNo":
						levelData.ItemConNo = intField.Value
					case "itemConsume":
						levelData.ItemConsume = intField.Value
					case "bulletConsume":
						levelData.BulletConsume = intField.Value
					case "bulletCount":
						levelData.BulletCount = intField.Value
					case "damage":
						levelData.Damage = intField.Value
					case "damagepc":
						levelData.DamagePC = intField.Value
					case "fixdamage":
						levelData.FixDamage = intField.Value
					case "criticalDamage":
						levelData.CriticalDamage = intField.Value
					case "attackCount":
						levelData.AttackCount = intField.Value
					case "mobCount":
						levelData.MobCount = intField.Value
					case "pad":
						levelData.PAD = intField.Value
					case "mad":
						levelData.MAD = intField.Value
					case "pdd":
						levelData.PDD = intField.Value
					case "mdd":
						levelData.MDD = intField.Value
					case "eva":
						levelData.EVA = intField.Value
					case "acc":
						levelData.ACC = intField.Value
					case "str":
						levelData.STR = intField.Value
					case "hp":
						levelData.HP = intField.Value
					case "mp":
						levelData.MP = intField.Value
					case "jump":
						levelData.Jump = intField.Value
					case "speed":
						levelData.Speed = intField.Value
					case "mastery":
						levelData.Mastery = intField.Value
					case "prop":
						levelData.Prop = intField.Value
					case "range":
						levelData.Range = intField.Value
					case "time":
						levelData.Time = time.Duration(intField.Value) * time.Second
					case "cooltime":
						levelData.Cooldown = time.Duration(intField.Value) * time.Second
					case "morph":
						levelData.Morph = intField.Value
					case "x":
						levelData.X = intField.Value
					case "y":
						levelData.Y = intField.Value
					case "z":
						levelData.Z = intField.Value
					}
				}

				for _, strField := range levelChild.Strings {
					switch strField.Name {
					case "damage":

						if val, err := strconv.Atoi(strField.Value); err == nil {
							levelData.Damage = val
						}
					case "attackCount":

						if val, err := strconv.Atoi(strField.Value); err == nil {
							levelData.AttackCount = val
						}
					case "acc":

						if val, err := strconv.Atoi(strField.Value); err == nil {
							levelData.ACC = val
						}
					case "time":

						if val, err := strconv.Atoi(strField.Value); err == nil {
							levelData.Time = time.Duration(val) * time.Second
						}
					case "hs":
						levelData.HS = strField.Value
					case "action":
						levelData.Action = strField.Value
					}
				}

				for _, vecField := range levelChild.Vectors {
					switch vecField.Name {
					case "lt":
						levelData.LT = types.Vector2[int32]{X: int32(vecField.X), Y: int32(vecField.Y)}
					case "rb":
						levelData.RB = types.Vector2[int32]{X: int32(vecField.X), Y: int32(vecField.Y)}
					}
				}

				skill.LevelData[levelNum] = levelData
			}
		} else {
			if skill.MaxLevel == 0 {
				skill.MaxLevel = 1
			}
		}

		skills[uint32(skillID)] = skill
	}

	return skills, nil
}

func getHardcodedExpTable() []uint32 {
	exp := []uint32{0, 15, 34, 57, 92, 135, 372, 560, 840, 1242, 1716, 2360, 3216, 4200, 5460, 7050, 8840, 11040, 13716, 16680, 20216, 24402, 28980, 34320, 40512, 47216, 54900, 63666, 73080, 83720, 95700, 108480, 122760, 138666, 155540, 174216, 194832, 216600, 240500, 266682, 294216, 324240, 356916, 391160, 428280, 468450, 510420, 555680, 604416, 655200, 709716, 748608, 789631, 832902, 878545, 926689, 977471, 1031036, 1087536, 1147032, 1209994, 1276301, 1346242, 1420016, 1497832, 1579913, 1666492, 1757815, 1854143, 1955750, 2062925, 2175973, 2295216, 2420993, 2553663, 2693603, 2841212, 2996910, 3161140, 3334370, 3517093, 3709829, 3913127, 4127566, 4353756, 4592341, 4844001, 5109452, 5389449, 5684790, 5996316, 6324914, 6671519, 7037118, 7422752, 7829518, 8258575, 8711144, 9188514, 9692044, 10223168, 10783397, 11374327, 11997640, 12655110, 13348610, 14080113, 14851703, 15665576, 16524049, 17429566, 18384706, 19392187, 20454878, 21575805, 22758159, 24005306, 25320796, 26708375, 28171993, 29715818, 31344244, 33061908, 34873700, 36784778, 38800583, 40926854, 43169645, 45535341, 48030677, 50662758, 53439077, 56367538, 59456479, 62714694, 66151459, 69776558, 73600313, 77633610, 81887931, 86375389, 91108760, 96101520, 101367883, 106992842, 112782213, 118962678, 125481832, 132358236, 139611467, 147262175, 155332142, 163844343, 172823012, 182293713, 192283408, 202820538, 213935103, 225658746, 238024845, 251068606, 264827165, 279339639, 294647508, 310794191, 327825712, 345790561, 364739883, 384727628, 405810702, 428049128, 451506220, 476248760, 502347192, 529875818, 558913012, 589541445, 621848316, 655925603, 691870326, 729784819, 769777027, 811960808, 856456260, 903390063, 952895838, 1005114529, 1060194805, 1118293480, 1179575962, 1244216724, 1312399800, 1384319309, 1460180007, 1540197871, 1624600714, 1713628833, 1807535693, 1906558648, 2011069705, 2121276324}

	expTable := make([]uint32, 201)
	copy(expTable, exp)
	if len(exp) < 201 {
		for i := len(exp); i < 201; i++ {
			expTable[i] = 0
		}
	}
	return expTable
}

func calculateDefaultExp(level int) uint32 {
	if level <= 0 || level > 200 {
		return 0
	}
	if level == 1 {
		return 15
	}

	baseExp := 15
	multiplier := 1.0
	if level > 10 {
		multiplier += float64(level-10) * 0.1
	}
	if level > 20 {
		multiplier += float64(level-20) * 0.15
	}
	if level > 30 {
		multiplier += float64(level-30) * 0.2
	}
	return uint32(float64(baseExp*level) * multiplier)
}
