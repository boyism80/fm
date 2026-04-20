package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type ShowItemGain struct {
	ItemId uint32
	Count  uint32
	Mode   constant.ShowItemGainType
}

func (p *ShowItemGain) Opcode() uint16 {
	switch p.Mode {
	case constant.ShowItemGainTypeChat:
		return 0x97
	case constant.ShowItemGainTypeStatus:
		return 0x1C
	default:
		return 0x1C
	}
}

func (p *ShowItemGain) Serialize(writer *stream.StreamWriter) error {
	switch p.Mode {
	case constant.ShowItemGainTypeChat:
		writer.WriteU8(3)
		writer.WriteU8(1)
		writer.WriteU32(p.ItemId)
		writer.WriteU32(p.Count)

	case constant.ShowItemGainTypeStatus:
		writer.WriteU16(0)
		writer.WriteU32(p.ItemId)
		writer.WriteU32(p.Count)
	}
	return nil
}

func (p *ShowItemGain) Deserialize(reader *stream.StreamReader) {
}
