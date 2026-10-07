package request

import (
	"github.com/boyism80/fm/stream"
)

type InspectCharacter struct {
	Tick        uint32
	CharacterID uint32
}

func (*InspectCharacter) Opcode() byte { return 0x50 }

func (p *InspectCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.WriteU32(p.CharacterID)
	return nil
}

func (p *InspectCharacter) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.CharacterID = reader.ReadU32()
}
