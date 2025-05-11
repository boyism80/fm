package data

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func loadConsumeFromXML(path string) (*[]*ConsumeSpec, error) {
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
		spec.Id = uint32(id)
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
