package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type ShowMesoGainType uint8

const (
	ShowMesoGainTypeStatus ShowMesoGainType = iota
	SHowMesoGainTypeChat
)

type ShowMesoGain struct {
	Count int32
	Mode  ShowMesoGainType
}

func (p *ShowMesoGain) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x1C)
	switch p.Mode {
	case SHowMesoGainTypeChat:
		writer.WriteU8(5)
		writer.Write32(p.Count)
		writer.Write32(-1)

	case ShowMesoGainTypeStatus:
		writer.WriteU8(0)
		writer.WriteU8(1)
		writer.Write32(p.Count)
		writer.WriteU16(0)
	}
	return nil
}

func (p *ShowMesoGain) Deserialize(reader *stream.StreamReader) error {
	return nil
}
