package response

import (
	"github.com/boyism80/fm/protocol/constant"
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
	w.WriteU8(uint8(constant.PartyS2CPartyPortal))
	w.WriteBoolean(!p.Animated)
	w.WriteU32(p.TownMapID)
	w.WriteU32(p.TargetMapID)
	w.WriteU32(p.SkillID)
	w.Write16(p.Position.X)
	w.Write16(p.Position.Y)
	return nil
}

func (p *PartyPortal) Deserialize(*stream.StreamReader) {}
