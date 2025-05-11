package data

import (
	"encoding/xml"
	"os"
	"strconv"
)

func loadStringNodeRecursive(root *node) []*StringSpec {

	templates := []*StringSpec{}
	for _, child := range root.Children {

		id, err := strconv.Atoi(child.Name)
		if err != nil {
			templates = append(templates, loadStringNodeRecursive(&child)...)
		}
		spec := StringSpec{
			Id: uint32(id),
		}
		for _, v := range child.Children {
			switch v.Name {
			case "name":
				spec.Name = v.Value
			case "desc":
				spec.Desc = v.Value
			}
		}
	}

	return templates
}

func loadStringFromXML(path string) (*[]*StringSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root node
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	templates := loadStringNodeRecursive(&root)
	return &templates, nil
}
