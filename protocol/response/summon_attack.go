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
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.SummonSkillID)
	w.WriteU8(p.Animation)
	w.WriteU8(uint8(len(p.Targets)))
	for _, t := range p.Targets {
		w.WriteU32(t.OID)
		w.WriteU8(6)
		w.WriteU32(t.Damage)
	}
	return nil
}

func (p *SummonAttack) Deserialize(*stream.StreamReader) {}
