package wz

import (
	"encoding/xml"
	"log"
	"os"
	"strconv"
)

type Commodity struct {
	SN          uint32
	ItemID      uint32
	Count       uint16
	Price       uint32
	Period      uint32
	Priority    uint8
	Gender      uint8
	OnSale      bool
	PayBackRate uint32
	Package     []uint32
}

type commodityNode struct {
	Ints []intField `xml:"int"`
}

type commodityRoot struct {
	XMLName xml.Name        `xml:"imgdir"`
	Entries []commodityNode `xml:"imgdir"`
}

type cashPackageNode struct {
	Name string          `xml:"name,attr"`
	SNs  []commodityNode `xml:"imgdir"`
}

type cashPackageRoot struct {
	XMLName  xml.Name          `xml:"imgdir"`
	Packages []cashPackageNode `xml:"imgdir"`
}

func loadCommodities(path string, packagePath string) (map[uint32]*Commodity, error) {
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
	byItem := make(map[uint32][]*Commodity)
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
			case "PbPoint":
				commodity.PayBackRate = uint32(field.Value)
			}
		}
		result[commodity.SN] = commodity
		byItem[commodity.ItemID] = append(byItem[commodity.ItemID], commodity)
	}

	packageFile, err := os.Open(packagePath)
	if err != nil {
		log.Printf("Failed to load CashPackage.img.xml: %v", err)
		return result, nil
	}
	defer packageFile.Close()

	var packages cashPackageRoot
	if err := xml.NewDecoder(packageFile).Decode(&packages); err != nil {
		log.Printf("Failed to load CashPackage.img.xml: %v", err)
		return result, nil
	}
	for _, node := range packages.Packages {
		itemID, err := strconv.ParseUint(node.Name, 10, 32)
		if err != nil {
			continue
		}
		var sns []uint32
		for _, list := range node.SNs {
			for _, field := range list.Ints {
				sns = append(sns, uint32(field.Value))
			}
		}
		for _, commodity := range byItem[uint32(itemID)] {
			commodity.Package = sns
		}
	}
	return result, nil
}
