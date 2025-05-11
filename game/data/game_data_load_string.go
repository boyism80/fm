package data

import (
	"encoding/xml"
	"os"
	"strconv"
)

func loadStringNodeRecursive(root *XMLNode) []*StringTemplate {

	templates := []*StringTemplate{}
	for _, child := range root.Children {

		id, err := strconv.Atoi(child.Name)
		if err != nil {
			templates = append(templates, loadStringNodeRecursive(&child)...)
		}
		template := StringTemplate{
			Id: uint32(id),
		}
		for _, v := range child.Children {
			switch v.Name {
			case "name":
				template.Name = v.Value
			case "desc":
				template.Desc = v.Value
			}
		}
	}

	return templates
}

func loadStringFromXML(path string) (*[]*StringTemplate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root XMLNode
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	templates := loadStringNodeRecursive(&root)
	return &templates, nil
}
