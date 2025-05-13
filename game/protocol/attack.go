package protocol

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type DamagePair struct {
	Damage  uint32
	Unknown bool
}

type AttackPair struct {
	ObjectId uint32
	Point    types.Vec2[int16]
	Attack   []DamagePair
}

type AttackInfo struct {
	Targets   uint8
	Hits      uint8
	Skill     uint32
	Charge    uint32
	Unk       uint8
	Speed     uint8
	Display   uint8
	LastTick  uint32
	Slot      uint8 // only for RANGED_ATTACK
	CsStar    uint8 // only for RANGED_ATTACK
	AOE       uint8 // only for RANGED_ATTACK
	AllDamage []AttackPair
	Position  types.Vec2[int16]
}

func (a *AttackInfo) Deserialize(sr *stream.StreamReader, opcode uint16) error {
	_ = sr.Skip(1) // unknown byte

	tbyte, err := sr.ReadU8()
	if err != nil {
		return err
	}
	a.Targets = (tbyte >> 4) & 0xF
	a.Hits = tbyte & 0xF

	skill, err := sr.ReadU32()
	if err != nil {
		return err
	}

	a.Skill = skill

	switch skill {
	case 2121001:
	case 2221001:
	case 2321001:
	case 3221001:
	case 3121004:
	case 13111002:
	case 5101004:
	case 15101003:
	case 5221004:
	case 5201002:
		charge, err := sr.ReadU32()
		if err != nil {
			return err
		}
		a.Charge = charge

	default:
		a.Charge = 0
	}

	err = sr.Skip(1) // nOption
	if err != nil {
		return err
	}

	unk, err := sr.ReadU8()
	if err != nil {
		return err
	}
	a.Unk = unk
	speed, err := sr.ReadU8()
	if err != nil {
		return err
	}
	a.Speed = speed
	display, err := sr.ReadU8()
	if err != nil {
		return err
	}
	a.Display = display
	lastTick, err := sr.ReadU32()
	if err != nil {
		return err
	}
	a.LastTick = lastTick

	a.AllDamage = make([]AttackPair, 0, a.Targets)
	for range int(a.Targets) {
		oid, err := sr.ReadU32()
		if err != nil {
			return err
		}
		if err := sr.Skip(14); err != nil {
			return err
		}

		attack := []DamagePair{}
		for range int(a.Hits) {
			damage, err := sr.ReadU32()
			if err != nil {
				return err
			}
			attack = append(attack, DamagePair{
				Damage:  damage,
				Unknown: false,
			})
		}

		a.AllDamage = append(a.AllDamage, AttackPair{
			ObjectId: oid,
			Attack:   attack,
		})
	}

	if opcode == 0x14 { // RANGED_ATTACK
		slot, err := sr.ReadU16()
		if err != nil {
			return err
		}
		a.Slot = uint8(slot)
		csStar, err := sr.ReadU16()
		if err != nil {
			return err
		}
		a.CsStar = uint8(csStar)
		aoe, err := sr.ReadU8()
		if err != nil {
			return err
		}
		a.AOE = aoe
	}

	if sr.Remaining() >= 4 {
		x, err := sr.Read16()
		if err != nil {
			return err
		}
		y, err := sr.Read16()
		if err != nil {
			return err
		}
		a.Position = types.Vec2[int16]{X: x, Y: y}
	}

	return nil
}
