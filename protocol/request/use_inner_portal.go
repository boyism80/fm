package request

import (
	"github.com/boyism80/fm/stream"
)

type UseInnerPortal struct {
	Mode       uint8
	PortalName string
	ToX        int16
	ToY        int16
	FromX      int16
	FromY      int16
}

func (p *UseInnerPortal) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseInnerPortal) Deserialize(reader *stream.StreamReader) {
	p.Mode = reader.ReadU8()
	p.PortalName = reader.ReadStr16()
	p.ToX = reader.Read16()
	p.ToY = reader.Read16()

	if reader.Remaining() >= 4 {
		p.FromX = reader.Read16()
		p.FromY = reader.Read16()
	}
}
