package request

import (
	"github.com/boyism80/fm/stream"
)

type DirectWarp struct {
	FieldKey   uint8
	PortalName string
}

func (*DirectWarp) Opcode() byte { return 0x53 }

func (p *DirectWarp) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.FieldKey)
	writer.WriteStr16(p.PortalName)
	return nil
}

func (p *DirectWarp) Deserialize(reader *stream.StreamReader) {
	p.FieldKey = reader.ReadU8()
	p.PortalName = reader.ReadStr16()
}
