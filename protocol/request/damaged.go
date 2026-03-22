package request

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

type Damaged struct {
	UpdateTick uint32
	Type       constant.IncomingHitType
	Element    constant.HitElement
	Damage     int32
	MobID      uint32
	OID        uint32
	Direction  uint8
	Reflect    uint8
	Level      uint8
	SkillID    uint8
}

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
	p.Type = constant.IncomingHitType(t)
	var elem uint8
	if elem, err = reader.ReadU8(); err != nil {
		return err
	}
	p.Element = constant.HitElement(elem)
	if p.Damage, err = reader.Read32(); err != nil {
		return err
	}

	switch p.Type {
	case constant.IncomingHitMapDebuff:
		if p.Level, err = reader.ReadU8(); err != nil {
			return err
		}
		if p.SkillID, err = reader.ReadU8(); err != nil {
			return err
		}

	case constant.IncomingHitEnv, constant.IncomingHitMist:
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
