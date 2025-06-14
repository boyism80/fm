package req

import (
	"github.com/boyism80/fm/common/stream"
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
	DamageTypeMist      DamageType = -4
	DamageTypeEnv       DamageType = -3
	DamageTypeMapDebuff DamageType = -2
	DamageTypeCollide   DamageType = -1
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
	case DamageTypeMapDebuff:
		if p.Level, err = reader.ReadU8(); err != nil {
			return err
		}
		if p.SkillID, err = reader.ReadU8(); err != nil {
			return err
		}

	case DamageTypeEnv, DamageTypeMist:
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
