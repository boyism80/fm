package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowRewardItemAnimation struct {
	CharacterID uint32
	ItemID      uint32
	Effect      string
}

func (p *ShowRewardItemAnimation) Opcode() uint16 {
	return 0x8F
}

func (p *ShowRewardItemAnimation) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(0x11)
	writer.WriteU32(p.ItemID)
	if p.Effect == "" {
		writer.WriteU8(0)
	} else {
		writer.WriteU8(1)
		writer.WriteStr16(p.Effect)
	}
	return nil
}

func (p *ShowRewardItemAnimation) Deserialize(reader *stream.StreamReader) {
}
