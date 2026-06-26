package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowNpcEffect struct {
	OID    uint32
	Action string
}

func (p *ShowNpcEffect) Opcode() uint16 {
	return 0xC0
}

func (p *ShowNpcEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.WriteStr16(p.Action)
	return nil
}

func (p *ShowNpcEffect) Deserialize(reader *stream.StreamReader) {
}
