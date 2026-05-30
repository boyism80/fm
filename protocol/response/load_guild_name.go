package response

import "github.com/boyism80/fm/stream"

type LoadGuildName struct {
	CharacterID uint32
	GuildName   string
}

func (p *LoadGuildName) Opcode() uint16 {
	return 0x93
}

func (p *LoadGuildName) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	if p.GuildName == "" {
		w.WriteU16(0)
		return nil
	}
	w.WriteStr16(p.GuildName)
	return nil
}

func (p *LoadGuildName) Deserialize(*stream.StreamReader) {}
