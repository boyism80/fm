package data

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func loadWeaponFromXML(path string) (*ItemTemplate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root XMLNode
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	template := ItemTemplate{}
	id, err := strconv.Atoi(strings.TrimSuffix(root.Name, ".img"))
	if err != nil {
		return nil, err
	}
	template.Id = uint32(id)
	node := root.Find("info")
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
			break

		case "reqJob":
			class, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Required.Class = class
		case "reqLevel":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Required.Class = value
		case "reqSTR":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Required.Str = value
		case "reqDEX":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Required.Dex = value
		case "reqINT":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Required.Int = value
		case "reqLUK":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Required.Luk = value
		case "incSTR":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.Str = value
		case "incDEX":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.Dex = value
		case "incINT":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.Int = value
		case "incLUK":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.Luk = value
		case "incPAD":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.PAD = value
		case "incMAD":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.MAD = value
		case "incPDD":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.PDD = value
		case "incMDD":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.MDD = value
		case "incPVPDamage":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.PVPDamage = value
		case "incSpeed":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.Speed = value
		case "incJump":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.Jump = value

		case "incACC":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.ACC = value

		case "incEVA":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.EVA = value
		case "incMHP":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.MaxHP = value
		case "incMMP":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Ability.MaxMP = value
		case "tuc":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.TUC = value
		case "price":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Price = value
		case "attackSpeed":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.AttackSpeed = value
		case "cash":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Cash = value == 1

		case "slotMax":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.SlotMax = uint16(value)

		case "quest":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Quest = value == 1

		case "equipTradeBlock":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.EquipTradeBlock = value == 1
		case "tradeAvailable":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.TradeAvailable = value
		case "hide":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.Hide = value == 1
		case "royalSpecial":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.RoyalSpecial = value == 1
		case "masterSpecial":
			value, err := strconv.Atoi(v.Value)
			if err != nil {
				return nil, err
			}
			template.MasterSpecial = value == 1
		default:
			log.Printf("%s is not declared in %s\n", v.Name, filepath.Base(path))
		}
	}

	return &template, nil
}
