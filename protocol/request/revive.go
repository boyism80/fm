package request

import (
	"github.com/boyism80/fm/stream"
)

type Revive struct {
	UnknownFlag uint8
	TargetID    int32
	PortalName  string
	SkipByte    uint8
	UseWheel    bool
}

func (p *Revive) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Revive) Deserialize(reader *stream.StreamReader) {
	p.UnknownFlag = reader.ReadU8()
	p.TargetID = reader.Read32()
	p.PortalName = reader.ReadStr16()
	reader.Skip(1)
	p.UseWheel = reader.ReadU8() > 0
}
