package wz

import (
	"encoding/xml"
	"os"
)

type Commodity struct {
	SN       uint32
	ItemID   uint32
	Count    uint16
	Price    uint32
	Period   uint32
	Priority uint8
	Gender   uint8
	OnSale   bool
}

type commodityNode struct {
	Ints []intField `xml:"int"`
}

type commodityRoot struct {
	XMLName xml.Name        `xml:"imgdir"`
	Entries []commodityNode `xml:"imgdir"`
}

func loadCommodities(path string) (map[uint32]*Commodity, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var root commodityRoot
	if err := xml.NewDecoder(file).Decode(&root); err != nil {
		return nil, err
	}

	result := make(map[uint32]*Commodity, len(root.Entries))
	for _, node := range root.Entries {
		commodity := &Commodity{}
		for _, field := range node.Ints {
			switch field.Name {
			case "SN":
				commodity.SN = uint32(field.Value)
			case "ItemId":
				commodity.ItemID = uint32(field.Value)
			case "Count":
				commodity.Count = uint16(field.Value)
			case "Price":
				commodity.Price = uint32(field.Value)
			case "Period":
				commodity.Period = uint32(field.Value)
			case "Priority":
				commodity.Priority = uint8(field.Value)
			case "Gender":
				commodity.Gender = uint8(field.Value)
			case "OnSale":
				commodity.OnSale = field.Value == 1
			}
		}
		result[commodity.SN] = commodity
	}
	return result, nil
}
