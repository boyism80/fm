package response

import (
	"github.com/boyism80/fm/stream"
)

// ShowMagnet is SHOW_MAGNET (0xB7), e.g. Monster Magnet pull result.
type ShowMagnet struct {
	MobID   uint32
	Success uint8
}

func (p *ShowMagnet) Opcode() uint16 {
	return 0xB7
}

func (p *ShowMagnet) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.MobID)
	writer.WriteU8(p.Success)
	return nil
}

func (p *ShowMagnet) Deserialize(reader *stream.StreamReader) {
}
