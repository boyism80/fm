package data

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func loadSpecialItemFromXML(path string) (*[]*SpecialItemTemplate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root XMLNode
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	templates := []*SpecialItemTemplate{}
	for _, v := range root.Children {
		template := SpecialItemTemplate{
			baseItemTemplate: &baseItemTemplate{
				SlotMax: 1,
			},
		}

		id, err := strconv.Atoi(v.Name)
		if err != nil {
			return nil, err
		}
		template.Id = uint32(id)
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

		templates = append(templates, &template)
	}

	return &templates, nil
}
