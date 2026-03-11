package response

import (
	"github.com/boyism80/fm/stream"
)

type CancelDebuff struct {
	OID    uint32
	Status int32
	Size   byte
}

func (p *CancelDebuff) Opcode() uint16 {
	return 0xB0
}

func (p *CancelDebuff) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.Write32(p.Status)
	writer.WriteU8(p.Size)
	return nil
}

func (p *CancelDebuff) Deserialize(reader *stream.StreamReader) error {
	return nil
}
