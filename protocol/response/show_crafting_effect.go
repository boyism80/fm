package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowCraftingEffect struct {
	CharacterID uint32
	Effect      string
	Time        int32
	Mode        int32
}

func (p *ShowCraftingEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowCraftingEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(0x1E)
	writer.WriteStr16(p.Effect)
	writer.Write32(p.Time)
	writer.Write32(p.Mode)
	return nil
}

func (p *ShowCraftingEffect) Deserialize(reader *stream.StreamReader) {
}
