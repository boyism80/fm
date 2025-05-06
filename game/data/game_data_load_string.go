package data

import (
	"encoding/xml"
	"os"
	"strconv"
)

func loadStringFromXML(path string) (*StringTemplate, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root XMLNode
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	m := StringTemplate{}
	for _, child := range root.Children {
		id, err := strconv.Atoi(child.Name)
		m.Id = uint32(id)
		if err != nil {
			return nil, err
		}
		for _, node := range child.Children {
			switch node.Name {
			case "name":
				m.Name = node.Value

			case "desc":
				m.Desc = node.Value
			}
		}
	}

	return &m, nil
}
