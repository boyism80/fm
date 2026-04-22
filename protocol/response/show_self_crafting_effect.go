package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowSelfCraftingEffect struct {
	Effect string
	Time   int32
	Mode   int32
}

func (p *ShowSelfCraftingEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfCraftingEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0x1E)
	writer.WriteStr16(p.Effect)
	writer.Write32(p.Time)
	writer.Write32(p.Mode)
	return nil
}

func (p *ShowSelfCraftingEffect) Deserialize(reader *stream.StreamReader) {
}
