package response

import "github.com/boyism80/fm/stream"

type SummonSkill struct {
	CharacterID uint32
	SummonOID   uint32
	NewStance   uint8
}

func (p *SummonSkill) Opcode() uint16 {
	return 0x7F
}

func (p *SummonSkill) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.SummonOID)
	w.WriteU8(p.NewStance)
	return nil
}

func (p *SummonSkill) Deserialize(*stream.StreamReader) {}
