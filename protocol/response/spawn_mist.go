package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type SpawnMist struct {
	OID        uint32
	PoisonMist uint8
	MobMist    bool
	CauserID   uint32
	SkillID    uint32
	SkillLevel uint8
	SkillDelay uint16
	Bounds     types.Rect[int32]
	MobSkill   bool
}

func (p *SpawnMist) Opcode() uint16 {
	return 0xCB
}

func (p *SpawnMist) Serialize(w *stream.StreamWriter) error {
	if err := w.WriteU32(p.OID); err != nil {
		return err
	}
	if err := w.WriteU8(p.PoisonMist); err != nil {
		return err
	}
	third := uint32(1)
	if p.PoisonMist == 2 || p.MobMist {
		third = p.CauserID
	}
	if err := w.WriteU32(third); err != nil {
		return err
	}
	if err := w.WriteU32(p.SkillID); err != nil {
		return err
	}
	if err := w.WriteU8(p.SkillLevel); err != nil {
		return err
	}
	if err := w.WriteU16(p.SkillDelay); err != nil {
		return err
	}
	if err := w.Write32(p.Bounds.Left); err != nil {
		return err
	}
	if err := w.Write32(p.Bounds.Top); err != nil {
		return err
	}
	if err := w.Write32(p.Bounds.Right); err != nil {
		return err
	}
	if err := w.Write32(p.Bounds.Bottom); err != nil {
		return err
	}
	tail := p.PoisonMist
	if p.MobMist {
		tail = 0
	}
	if err := w.WriteU8(tail); err != nil {
		return err
	}
	if !p.MobSkill {
		smoke := uint8(0)
		if p.SkillID == 4221006 {
			smoke = 2
		}
		if err := w.WriteU8(smoke); err != nil {
			return err
		}
		if err := w.WriteU32(0); err != nil {
			return err
		}
	} else {
		if err := w.WriteU64(0); err != nil {
			return err
		}
	}
	if err := w.WriteU32(p.CauserID); err != nil {
		return err
	}
	return nil
}

func (p *SpawnMist) Deserialize(*stream.StreamReader) {}
