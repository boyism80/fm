package response

import (
	"github.com/boyism80/fm/stream"
)

type FieldRelocate struct {
	Portal uint8
}

func (p *FieldRelocate) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU8(p.Portal)
	return nil
}

func (p *FieldRelocate) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Portal = reader.ReadU8()
}

func (p *FieldRelocate) Opcode() uint16 {
	return 0x98
}
