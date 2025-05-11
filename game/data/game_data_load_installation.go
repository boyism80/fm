package data

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func loadInstallationFromXML(path string) (*[]*InstallationSpec, error) {
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
		spec.Id = uint32(id)
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
