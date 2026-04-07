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

func (p *Damaged) Deserialize(reader *stream.StreamReader) {
	p.UpdateTick = reader.ReadU32()
	var t int8
	t = reader.Read8()
	p.Type = constant.IncomingHitType(t)
	var elem uint8
	elem = reader.ReadU8()
	p.Element = constant.HitElement(elem)
	p.Damage = reader.Read32()

	switch p.Type {
	case constant.IncomingHitMapDebuff:
		p.Level = reader.ReadU8()
		p.SkillID = reader.ReadU8()

	case constant.IncomingHitEnv, constant.IncomingHitMist:
	default:
		p.MobID = reader.ReadU32()
		p.OID = reader.ReadU32()
		p.Direction = reader.ReadU8()
		p.Reflect = reader.ReadU8()
	}
}
