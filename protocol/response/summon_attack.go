package response

import "github.com/boyism80/fm/stream"

type SummonAttackTarget struct {
	OID    uint32
	Damage uint32
}

type SummonAttack struct {
	CharacterID   uint32
	SummonSkillID uint32
	Animation     uint8
	Targets       []SummonAttackTarget
}

func (p *SummonAttack) Opcode() uint16 {
	return 0x7E
}

func (p *SummonAttack) Serialize(w *stream.StreamWriter) error {
	if err := w.WriteU32(p.CharacterID); err != nil {
		return err
	}
	if err := w.WriteU32(p.SummonSkillID); err != nil {
		return err
	}
	if err := w.WriteU8(p.Animation); err != nil {
		return err
	}
	if err := w.WriteU8(uint8(len(p.Targets))); err != nil {
		return err
	}
	for _, t := range p.Targets {
		if err := w.WriteU32(t.OID); err != nil {
			return err
		}
		if err := w.WriteU8(6); err != nil {
			return err
		}
		if err := w.WriteU32(t.Damage); err != nil {
			return err
		}
	}
	return nil
}

func (p *SummonAttack) Deserialize(*stream.StreamReader) {}
