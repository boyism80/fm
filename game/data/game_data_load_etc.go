package data

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func loadGeneralItemFromXML(path string) (*[]*GeneralItemSpec, error) {
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
		spec.Id = uint32(id)
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
