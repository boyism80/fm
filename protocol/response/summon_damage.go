package response

import "github.com/boyism80/fm/stream"

type DamageSummon struct {
	CharacterID   uint32
	SummonSkillID uint32
	Unknown       uint8
	Damage        uint32
	MonsterIDFrom uint32
}

func (p *DamageSummon) Opcode() uint16 {
	return 0x80
}

func (p *DamageSummon) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.CharacterID)
	w.WriteU32(p.SummonSkillID)
	w.WriteU8(p.Unknown)
	w.WriteU32(p.Damage)
	w.WriteU32(p.MonsterIDFrom)
	w.WriteU8(0)
	return nil
}

func (p *DamageSummon) Deserialize(*stream.StreamReader) {}
