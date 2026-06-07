package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type TriggerReactor struct {
	Reactor *dto.Reactor
	Stance  int32
}

func (p *TriggerReactor) Opcode() uint16 {
	return 0xCF
}

func (p *TriggerReactor) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Reactor.OID)
	writer.WriteU8(p.Reactor.State)
	writer.Write16(p.Reactor.Position.X)
	writer.Write16(p.Reactor.Position.Y)
	writer.Write32(p.Stance)
	return nil
}

func (p *TriggerReactor) Deserialize(*stream.StreamReader) {}
