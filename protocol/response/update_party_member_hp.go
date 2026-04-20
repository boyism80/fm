package response

import "github.com/boyism80/fm/stream"

type UpdatePartyMemberHP struct {
	CharacterID uint32
	CurrentHP   int32
	MaxHP       int32
}

func (p *UpdatePartyMemberHP) Opcode() uint16 {
	return 0x92
}

func (p *UpdatePartyMemberHP) Serialize(w *stream.StreamWriter) error {
	w.Write32(int32(p.CharacterID))
	w.Write32(p.CurrentHP)
	w.Write32(p.MaxHP)
	return nil
}

func (p *UpdatePartyMemberHP) Deserialize(*stream.StreamReader) {}
