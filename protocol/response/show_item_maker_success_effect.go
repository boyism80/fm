package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowItemMakerSuccessEffect struct {
	CharacterID uint32
}

func (p *ShowItemMakerSuccessEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowItemMakerSuccessEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(0x11)
	writer.WriteU32(0)
	return nil
}

func (p *ShowItemMakerSuccessEffect) Deserialize(reader *stream.StreamReader) {
}
