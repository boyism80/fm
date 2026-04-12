package dto

import (
	"github.com/boyism80/fm/stream"
)

type Item interface {
	Serialize(writer *stream.StreamWriter, trade bool, slot int16)
	GetCount() uint16
}
