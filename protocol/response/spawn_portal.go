package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

const DisabledPortalMapID uint32 = 999999999

type SpawnPortal struct {
	TownMapID   uint32
	TargetMapID uint32
	SkillID     uint32
	Position    *types.Vector2[int16]
}

func (p *SpawnPortal) Opcode() uint16 {
	return 0x32
}

func (p *SpawnPortal) Serialize(w *stream.StreamWriter) error {
	if err := w.WriteU32(p.TownMapID); err != nil {
		return err
	}
	if err := w.WriteU32(p.TargetMapID); err != nil {
		return err
	}
	if p.TownMapID != DisabledPortalMapID && p.TargetMapID != DisabledPortalMapID {
		if err := w.WriteU32(p.SkillID); err != nil {
			return err
		}
		if p.Position == nil {
			if err := w.Write16(0); err != nil {
				return err
			}
			return w.Write16(0)
		}
		if err := w.Write16(p.Position.X); err != nil {
			return err
		}
		if err := w.Write16(p.Position.Y); err != nil {
			return err
		}
	}
	return nil
}

func (p *SpawnPortal) Deserialize(*stream.StreamReader) {}
