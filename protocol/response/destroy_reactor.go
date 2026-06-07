package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type DestroyReactor struct {
	Reactor *dto.Reactor
}

func (p *DestroyReactor) Opcode() uint16 {
	return 0xD2
}

func (p *DestroyReactor) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Reactor.OID)
	writer.WriteU8(p.Reactor.State)
	writer.Write16(p.Reactor.Position.X)
	writer.Write16(p.Reactor.Position.Y)
	return nil
}

func (p *DestroyReactor) Deserialize(*stream.StreamReader) {}
