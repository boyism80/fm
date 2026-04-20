package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type ItemGainFailed struct {
	Mode constant.ItemGainFailedType
}

func (p *ItemGainFailed) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU8(uint8(p.Mode))
	writer.WriteU16(0)
	return nil
}

func (p *ItemGainFailed) Deserialize(reader *stream.StreamReader) {
}

func (p *ItemGainFailed) Opcode() uint16 {
	return 0x1C
}
