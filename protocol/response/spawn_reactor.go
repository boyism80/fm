package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type SpawnReactor struct {
	Reactor *dto.Reactor
}

func (p *SpawnReactor) Opcode() uint16 {
	return 0xD1
}

func (p *SpawnReactor) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Reactor.OID)
	writer.WriteU32(p.Reactor.ReactorID)
	writer.WriteU8(p.Reactor.State)
	writer.Write16(p.Reactor.Position.X)
	writer.Write16(p.Reactor.Position.Y)
	writer.WriteU8(p.Reactor.Facing)
	writer.WriteStr16(p.Reactor.Name)
	return nil
}

func (p *SpawnReactor) Deserialize(*stream.StreamReader) {}
