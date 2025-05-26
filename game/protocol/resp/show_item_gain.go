package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type ShowItemGainType uint8

const (
	ShowItemGainTypeStatus ShowItemGainType = iota
	SHowItemGainTypeChat
)

type ShowItemGain struct {
	ItemId uint32
	Count  uint32
	Mode   ShowItemGainType
}

func (p *ShowItemGain) Serialize(writer *stream.StreamWriter) error {
	switch p.Mode {
	case SHowItemGainTypeChat:
		writer.WriteU16(0x97)
		writer.WriteU8(3)
		writer.WriteU8(1)
		writer.WriteU32(p.ItemId)
		writer.WriteU32(p.Count)

	case ShowItemGainTypeStatus:
		writer.WriteU16(0x1C)
		writer.WriteU16(0)
		writer.WriteU32(p.ItemId)
		writer.WriteU32(p.Count)
	}
	return nil
}

func (p *ShowItemGain) Deserialize(reader *stream.StreamReader) error {
	return nil
}
