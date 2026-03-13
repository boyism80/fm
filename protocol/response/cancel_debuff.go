package response

import (
	"github.com/boyism80/fm/stream"
)

type CancelMobStatus struct {
	OID    uint32
	Status int32
	Size   byte
}

func (p *CancelMobStatus) Opcode() uint16 {
	return 0xB0
}

func (p *CancelMobStatus) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.Write32(p.Status)
	writer.WriteU8(p.Size)
	return nil
}

func (p *CancelMobStatus) Deserialize(reader *stream.StreamReader) error {
	return nil
}
