package response

import (
	"github.com/boyism80/fm/game/constant"
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
	if err := w.WriteU32(p.OwnerID); err != nil {
		return err
	}
	if err := w.WriteU32(p.OID); err != nil {
		return err
	}
	if err := w.WriteU32(uint32(p.SkillID)); err != nil {
		return err
	}
	if err := w.WriteU8(p.SkillLevel); err != nil {
		return err
	}
	if err := w.Write16(p.Position.X); err != nil {
		return err
	}
	if err := w.Write16(p.Position.Y); err != nil {
		return err
	}
	if err := w.WriteU8(0); err != nil {
		return err
	}
	if err := w.Write16(0); err != nil {
		return err
	}
	if err := w.WriteU8(uint8(p.MovementType)); err != nil {
		return err
	}
	if err := w.WriteU8(uint8(p.SummonType)); err != nil {
		return err
	}
	if p.Animated {
		if err := w.WriteU8(1); err != nil {
			return err
		}
	} else {
		if err := w.WriteU8(0); err != nil {
			return err
		}
		if err := w.Write(make([]byte, 16)); err != nil {
			return err
		}
	}
	return nil
}

func (p *SpawnSummon) Deserialize(*stream.StreamReader) error { return nil }
