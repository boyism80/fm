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
	if err := w.WriteU32(p.CharacterID); err != nil {
		return err
	}
	if err := w.WriteU32(p.SummonSkillID); err != nil {
		return err
	}
	if err := w.WriteU8(p.Unknown); err != nil {
		return err
	}
	if err := w.WriteU32(p.Damage); err != nil {
		return err
	}
	if err := w.WriteU32(p.MonsterIDFrom); err != nil {
		return err
	}
	if err := w.WriteU8(0); err != nil {
		return err
	}
	return nil
}

func (p *DamageSummon) Deserialize(*stream.StreamReader) error { return nil }
