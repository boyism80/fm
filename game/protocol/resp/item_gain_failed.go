package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type ItemGainFailedType uint8

const (
	ITEM_GAIN_FAILED_TYPE_FULL  ItemGainFailedType = 0xFF
	ITEM_GAIN_FAILED_TYPE_ERROR ItemGainFailedType = 0xFE
)

type ItemGainFailed struct {
	Mode ItemGainFailedType
}

func (p *ItemGainFailed) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU8(uint8(p.Mode))
	writer.WriteU16(0)
	return nil
}

func (p *ItemGainFailed) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *ItemGainFailed) Opcode() uint16 {
	return 0x1C
}
