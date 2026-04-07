package request

import (
	"github.com/boyism80/fm/stream"
)

type Warp struct {
	Reason     uint8
	Target     uint32
	PortalName string
	Wheel      bool
}

func (p *Warp) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Warp) Deserialize(reader *stream.StreamReader) {
	p.Reason = reader.ReadU8()
	p.Target = reader.ReadU32()
	p.PortalName = reader.ReadStr16()
	reader.Skip(1)
	p.Wheel = reader.ReadBool()
}
