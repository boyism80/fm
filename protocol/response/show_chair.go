package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowChair struct {
	CharacterID uint32
	ItemID      uint32
}

func (p *ShowChair) Opcode() uint16 {
	return 0x8D
}

func (p *ShowChair) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *ShowChair) Deserialize(reader *stream.StreamReader) error {
	return nil
}
