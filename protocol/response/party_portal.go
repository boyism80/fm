package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type PartyPortal struct {
	TownMapID   uint32
	TargetMapID uint32
	SkillID     uint32
	Position    types.Vector2[int16]
	Animated    bool
}

func (p *PartyPortal) Opcode() uint16 {
	return 0x2D
}

func (p *PartyPortal) Serialize(w *stream.StreamWriter) error {
	if err := w.WriteU8(34); err != nil {
		return err
	}
	if p.Animated {
		if err := w.WriteU8(0); err != nil {
			return err
		}
	} else {
		if err := w.WriteU8(1); err != nil {
			return err
		}
	}
	if err := w.WriteU32(p.TownMapID); err != nil {
		return err
	}
	if err := w.WriteU32(p.TargetMapID); err != nil {
		return err
	}
	if err := w.WriteU32(p.SkillID); err != nil {
		return err
	}
	if err := w.Write16(p.Position.X); err != nil {
		return err
	}
	return w.Write16(p.Position.Y)
}

func (p *PartyPortal) Deserialize(*stream.StreamReader) error { return nil }
