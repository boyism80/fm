package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type ItemGainFailedType uint8

const (
	ItemGainFailedTypeFull  ItemGainFailedType = 0xFF
	ItemGainFailedTypeError ItemGainFailedType = 0xFE
)

type ItemGainFailed struct {
	Mode ItemGainFailedType
}

func (p *ItemGainFailed) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x1C)
	writer.WriteU8(0)
	writer.WriteU8(uint8(p.Mode))
	writer.WriteU16(0)
	return nil
}

func (p *ItemGainFailed) Deserialize(reader *stream.StreamReader) error {
	return nil
}
