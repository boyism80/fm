package request

import (
	"github.com/boyism80/fm/stream"
)

type DirectWarp struct {
	Mode       uint8
	PortalName string
}

func (*DirectWarp) Opcode() byte { return 0x53 }

func (p *DirectWarp) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DirectWarp) Deserialize(reader *stream.StreamReader) {
	p.Mode = reader.ReadU8()
	p.PortalName = reader.ReadStr16()
}
