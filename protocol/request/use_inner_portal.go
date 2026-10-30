package request

import (
	"github.com/boyism80/fm/stream"
)

// UseInnerPortal is a portal to another spot of the same map. The client sends where the character stands, then
// where the target portal is, and moves itself there.
type UseInnerPortal struct {
	FieldKey   uint8
	PortalName string
	FromX      int16
	FromY      int16
	ToX        int16
	ToY        int16
}

func (*UseInnerPortal) Opcode() byte { return 0x54 }

func (p *UseInnerPortal) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.FieldKey)
	writer.WriteStr16(p.PortalName)
	writer.Write16(p.FromX)
	writer.Write16(p.FromY)
	writer.Write16(p.ToX)
	writer.Write16(p.ToY)
	return nil
}

func (p *UseInnerPortal) Deserialize(reader *stream.StreamReader) {
	p.FieldKey = reader.ReadU8()
	p.PortalName = reader.ReadStr16()
	p.FromX = reader.Read16()
	p.FromY = reader.Read16()
	p.ToX = reader.Read16()
	p.ToY = reader.Read16()
}
