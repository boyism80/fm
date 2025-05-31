package req

import (
	"github.com/boyism80/fm/common/stream"
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

func (p *Warp) Deserialize(reader *stream.StreamReader) error {
	p.Reason, _ = reader.ReadU8()
	p.Target, _ = reader.ReadU32()
	p.PortalName, _ = reader.ReadStr16()
	reader.Skip(1)
	p.Wheel, _ = reader.ReadBool()
	return nil
}
