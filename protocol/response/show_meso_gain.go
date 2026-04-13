package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type ShowMesoGain struct {
	Count int32
	Mode  constant.ShowMesoGainType
}

func (p *ShowMesoGain) Serialize(writer *stream.StreamWriter) error {
	switch p.Mode {
	case constant.ShowMesoGainTypeChat:
		writer.WriteU8(5)
		writer.Write32(p.Count)
		writer.Write32(-1)

	case constant.ShowMesoGainTypeStatus:
		writer.WriteU8(0)
		writer.WriteU8(1)
		writer.Write32(p.Count)
		writer.WriteU16(0)
	}
	return nil
}

func (p *ShowMesoGain) Deserialize(reader *stream.StreamReader) {
}

func (p *ShowMesoGain) Opcode() uint16 {
	return 0x1C
}
