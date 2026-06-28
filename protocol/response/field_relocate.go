package response

import (
	"github.com/boyism80/fm/stream"
)

type FieldRelocate struct {
	Portal uint8
}

func (p *FieldRelocate) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0)
	writer.WriteU32(uint32(p.Portal))
	return nil
}

func (p *FieldRelocate) Deserialize(reader *stream.StreamReader) {
}

func (p *FieldRelocate) Opcode() uint16 {
	return 0x56
}
