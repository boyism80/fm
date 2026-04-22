package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowSelfRewardItemAnimation struct {
	ItemID uint32
	Effect string
}

func (p *ShowSelfRewardItemAnimation) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfRewardItemAnimation) Serialize(writer *stream.StreamWriter) error {
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

func (p *ShowSelfRewardItemAnimation) Deserialize(reader *stream.StreamReader) {
}
