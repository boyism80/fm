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
	if err := w.WriteU32(p.CharacterID); err != nil {
		return err
	}
	if err := w.WriteU32(p.SummonOID); err != nil {
		return err
	}
	if err := w.WriteU8(p.NewStance); err != nil {
		return err
	}
	return nil
}

func (p *SummonSkill) Deserialize(*stream.StreamReader) {}
