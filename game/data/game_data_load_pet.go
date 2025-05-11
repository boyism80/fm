package data

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func loadPetFromXML(path string) (*PetTemplate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root XMLNode
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(strings.TrimSuffix(root.Name, ".img"))
	if err != nil {
		return nil, err
	}
	template := &PetTemplate{
		baseItemTemplate: &baseItemTemplate{
			Id: uint32(id),
		},
	}
	info := root.Find("info")
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
	return template, nil
}
