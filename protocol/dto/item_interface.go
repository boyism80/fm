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
}
