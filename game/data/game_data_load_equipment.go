package data

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var mutex sync.Mutex = sync.Mutex{}
var visit map[string]bool = map[string]bool{}

func loadWeaponFromXML(path string) (*EquipmentSpec, error) {
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
	spec.Id = uint32(id)
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
