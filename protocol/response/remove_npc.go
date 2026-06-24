package response

import (
	"github.com/boyism80/fm/stream"
)

type RemoveNpc struct {
	OID uint32
}

func (p *RemoveNpc) Opcode() uint16 {
	return 0xBC
}

func (p *RemoveNpc) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	return nil
}

func (s *RemoveNpc) Deserialize(reader *stream.StreamReader) {
}
