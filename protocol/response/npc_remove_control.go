package response

import (
	"github.com/boyism80/fm/stream"
)

type NpcRemoveControl struct {
	OID uint32
}

func (p *NpcRemoveControl) Opcode() uint16 {
	return 0xBD
}

func (p *NpcRemoveControl) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU32(p.OID)
	return nil
}

func (s *NpcRemoveControl) Deserialize(reader *stream.StreamReader) {
}
