package response

import "github.com/boyism80/fm/stream"

type ShowScrollEffect struct {
	CharacterID      uint32
	Success          bool
	DestroyedByCurse bool
	LegendarySpirit  bool
	WhiteScroll      bool
}

func (p *ShowScrollEffect) Opcode() uint16 {
	return 0x74
}

func (p *ShowScrollEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteBoolean(p.Success)
	writer.WriteBoolean(p.DestroyedByCurse)
	writer.WriteBoolean(p.LegendarySpirit)
	writer.WriteBoolean(p.WhiteScroll)
	writer.WriteU32(0)
	return nil
}

func (p *ShowScrollEffect) Deserialize(reader *stream.StreamReader) {
}
