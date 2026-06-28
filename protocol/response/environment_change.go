package response

import (
	"github.com/boyism80/fm/stream"
)

type EnvironmentChangeMode uint8

const (
	EnvironmentChangeModeObjectState EnvironmentChangeMode = 2
	EnvironmentChangeModeMapEffect   EnvironmentChangeMode = 3
	EnvironmentChangeModeSound       EnvironmentChangeMode = 4
	EnvironmentChangeModeMusic       EnvironmentChangeMode = 6
)

type EnvironmentChange struct {
	Mode EnvironmentChangeMode
	Env  string
}

func (p *EnvironmentChange) Opcode() uint16 {
	return 0x5F
}

func (p *EnvironmentChange) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	writer.WriteStr16(p.Env)
	return nil
}

func (p *EnvironmentChange) Deserialize(*stream.StreamReader) {}
