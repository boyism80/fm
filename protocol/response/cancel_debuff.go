package response

import (
	"github.com/boyism80/fm/stream"
)

type CancelMobBuff struct {
	OID    uint32
	Status int32
	Size   byte
}

func (p *CancelMobBuff) Opcode() uint16 {
	return 0xB0
}

func (p *CancelMobBuff) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.OID)
	writer.Write32(p.Status)
	writer.WriteU8(p.Size)
	return nil
}

func (p *CancelMobBuff) Deserialize(reader *stream.StreamReader) {
}
