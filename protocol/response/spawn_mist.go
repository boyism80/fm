package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type SpawnMist struct {
	OID        uint32
	Type       constant.MistType
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
	w.WriteU32(p.OID)
	w.WriteU8(uint8(p.Type))
	if p.Type == constant.MistTypeSmoke || p.MobMist {
		w.WriteU32(p.CauserID)
	} else {
		w.WriteU32(1)
	}
	w.WriteU32(p.SkillID)
	w.WriteU8(p.SkillLevel)
	w.WriteU16(p.SkillDelay)
	w.Write32(p.Bounds.Left)
	w.Write32(p.Bounds.Top)
	w.Write32(p.Bounds.Right)
	w.Write32(p.Bounds.Bottom)
	if p.MobMist {
		w.WriteU8(0)
	} else {
		w.WriteU8(uint8(p.Type))
	}
	if !p.MobSkill {
		w.WriteU8(uint8(p.Type))
		w.WriteU32(0)
	} else {
		w.WriteU64(0)
	}
	w.WriteU32(p.CauserID)
	return nil
}

func (p *SpawnMist) Deserialize(*stream.StreamReader) {}
