package response

import "github.com/boyism80/fm/stream"

type SummonAttackTarget struct {
	OID    uint32
	Damage uint32
}

type SummonAttack struct {
	CharacterID uint32
	SummonOID   uint32
	Animation   uint8
	Targets     []SummonAttackTarget
}

func (p *SummonAttack) Opcode() uint16 {
	return 0x7E
}

func (p *SummonAttack) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.SummonOID)
	w.WriteU8(p.Animation)
	w.WriteU8(uint8(len(p.Targets)))
	for _, t := range p.Targets {
		w.WriteU32(t.OID)
		w.WriteU8(6)
		w.WriteU32(t.Damage)
	}
	return nil
}

func (p *SummonAttack) Deserialize(reader *stream.StreamReader) {
	p.CharacterID = reader.ReadU32()
	p.SummonOID = reader.ReadU32()
	p.Animation = reader.ReadU8()
	p.Targets = make([]SummonAttackTarget, reader.ReadU8())
	for i := range p.Targets {
		p.Targets[i].OID = reader.ReadU32()
		reader.ReadU8()
		p.Targets[i].Damage = reader.ReadU32()
	}
}
