package dto

import (
	"github.com/boyism80/fm/stream"
)

type SlotEncodeMode uint8

const (
	SlotEncodeActual SlotEncodeMode = iota
	SlotEncodeZero
	SlotEncodeOmit
)

type ItemSerializeOption struct {
	Trade    bool
	Slot     int16
	SlotMode SlotEncodeMode
}

type Item interface {
	Serialize(writer *stream.StreamWriter, opt ItemSerializeOption)
	GetCount() uint16
	GetItemID() uint32
}

func NewItemFromStream(reader *stream.StreamReader) Item {
	kind := reader.ReadU8()
	itemID := reader.ReadU32()
	switch {
	case kind == 1:
		item := &Equipment{ItemId: itemID}
		item.Deserialize(reader)
		return item
	case kind == 3:
		item := &PetItem{ItemId: itemID}
		item.Deserialize(reader)
		return item
	case itemID/1000000 == 2:
		item := &ConsumeItem{ItemId: itemID}
		item.Deserialize(reader)
		return item
	case itemID/1000000 == 3:
		item := &InstallationItem{ItemId: itemID}
		item.Deserialize(reader)
		return item
	case itemID/1000000 == 5:
		item := &CashItem{ItemId: itemID}
		item.Deserialize(reader)
		return item
	default:
		item := &MiscItem{ItemId: itemID}
		item.Deserialize(reader)
		return item
	}
}
