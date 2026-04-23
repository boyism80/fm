package request

import (
	"github.com/boyism80/fm/stream"
)

type DropMeso struct {
	Tick  uint32
	Count int32
}

func (*DropMeso) Opcode() byte { return 0x4D }

func (p *DropMeso) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DropMeso) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Count = reader.Read32()
}
