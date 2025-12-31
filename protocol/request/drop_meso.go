package request

import (
	"github.com/boyism80/fm/stream"
)

type DropMeso struct {
	Tick  uint32
	Count int32
}

func (p *DropMeso) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DropMeso) Deserialize(reader *stream.StreamReader) error {
	p.Tick, _ = reader.ReadU32()
	p.Count, _ = reader.Read32()
	return nil
}
