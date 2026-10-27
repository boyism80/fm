package request

import (
	"github.com/boyism80/fm/stream"
)

type Aggro struct {
	OID uint32
}

func (*Aggro) Opcode() byte { return 0x96 }

func (p *Aggro) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	return nil
}

func (p *Aggro) Deserialize(reader *stream.StreamReader) {
	p.OID = reader.ReadU32()
}
