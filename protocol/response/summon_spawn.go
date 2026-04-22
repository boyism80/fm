package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type SpawnSummon struct {
	OwnerID      uint32
	OID          uint32
	SkillID      constant.SkillID
	SkillLevel   uint8
	Position     types.Vector2[int16]
	MovementType constant.SummonMovementType
	SummonType   constant.SummonType
	Animated     bool
}

func (p *SpawnSummon) Opcode() uint16 {
	return 0x7B
}

func (p *SpawnSummon) Serialize(w *stream.StreamWriter) error {
	w.WriteU32(p.OwnerID)
	w.WriteU32(p.OID)
	w.WriteU32(uint32(p.SkillID))
	w.WriteU8(p.SkillLevel)
	w.Write16(p.Position.X)
	w.Write16(p.Position.Y)
	w.WriteU8(0)
	w.Write16(0)
	w.WriteU8(uint8(p.MovementType))
	w.WriteU8(uint8(p.SummonType))
	w.WriteBoolean(p.Animated)
	if !p.Animated {
		w.Write(make([]byte, 16))
	}
	return nil
}

func (p *SpawnSummon) Deserialize(*stream.StreamReader) {}
