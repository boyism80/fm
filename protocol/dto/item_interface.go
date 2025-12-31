package dto

import (
	"github.com/boyism80/fm/stream"
)

// Item represents an item DTO interface
type Item interface {
	Serialize(writer *stream.StreamWriter, trade bool, slot int16)
	GetCount() uint16
}
