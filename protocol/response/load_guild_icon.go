package response

import "github.com/boyism80/fm/stream"

type LoadGuildIcon struct {
	CharacterID  uint32
	LogoBG       uint16
	LogoBGColor  uint8
	Logo         uint16
	LogoColor    uint8
	HasGuildIcon bool
}

func (p *LoadGuildIcon) Opcode() uint16 {
	return 0x94
}

func (p *LoadGuildIcon) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	if !p.HasGuildIcon {
		w.Write(make([]byte, 6))
		return nil
	}
	w.WriteU16(p.LogoBG)
	w.WriteU8(p.LogoBGColor)
	w.WriteU16(p.Logo)
	w.WriteU8(p.LogoColor)
	return nil
}

func (p *LoadGuildIcon) Deserialize(*stream.StreamReader) {}
