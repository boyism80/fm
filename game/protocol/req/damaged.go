package req

import (
	"github.com/boyism80/fm/core/stream"
)

type DamageType int8

type Damaged struct {
	UpdateTick uint32
	Type       DamageType
	Element    uint8
	Damage     int32

	MobID     uint32
	OID       uint32
	Direction uint8
	Reflect   uint8

	Level   uint8
	SkillID uint8
}

const (
	DAMAGE_TYPE_MIST       DamageType = -4
	DAMAGE_TYPE_ENV        DamageType = -3
	DAMAGE_TYPE_MAP_DEBUFF DamageType = -2
	DAMAGE_TYPE_COLLIDE    DamageType = -1
)

func (p *Damaged) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *Damaged) Deserialize(reader *stream.StreamReader) error {
	var err error
	if p.UpdateTick, err = reader.ReadU32(); err != nil {
		return err
	}
	var t int8
	if t, err = reader.Read8(); err != nil {
		return err
	}
	p.Type = DamageType(t)
	if p.Element, err = reader.ReadU8(); err != nil {
		return err
	}
	if p.Damage, err = reader.Read32(); err != nil {
		return err
	}

	switch p.Type {
	case DAMAGE_TYPE_MAP_DEBUFF:
		if p.Level, err = reader.ReadU8(); err != nil {
			return err
		}
		if p.SkillID, err = reader.ReadU8(); err != nil {
			return err
		}

	case DAMAGE_TYPE_ENV, DAMAGE_TYPE_MIST:
	default:
		if p.MobID, err = reader.ReadU32(); err != nil {
			return err
		}
		if p.OID, err = reader.ReadU32(); err != nil {
			return err
		}
		if p.Direction, err = reader.ReadU8(); err != nil {
			return err
		}
		if p.Reflect, err = reader.ReadU8(); err != nil {
			return err
		}
	}
	return nil
}
