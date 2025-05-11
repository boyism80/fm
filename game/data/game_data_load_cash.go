package data

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func loadCashItemFromXML(path string) (*[]*CashItemTemplate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root XMLNode
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	templates := []*CashItemTemplate{}
	for _, v := range root.Children {
		template := CashItemTemplate{
			baseItemTemplate: &baseItemTemplate{},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		template.Id = uint32(id)
		info := v.Find("info")
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
				template.SlotMax = uint16(slotMax)

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

		templates = append(templates, &template)
	}

	return &templates, nil
}
