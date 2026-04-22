package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowSelfEffect struct {
	Type EffectType
}

func (p *ShowSelfEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Type))
	return nil
}

func (p *ShowSelfEffect) Deserialize(reader *stream.StreamReader) {
}
