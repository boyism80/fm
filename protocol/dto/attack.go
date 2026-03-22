package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type DamagePair struct {
	Damage  uint32
	Unknown bool
}

type AttackPair struct {
	OID         uint32
	DamagePairs []DamagePair
}

type AttackInfo struct {
	Targets  uint8
	Hits     uint8
	Skill    uint32
	Charge   uint32
	Unk      uint8
	Speed    uint8
	Display  uint8
	LastTick uint32
	Slot     uint16
	CsStar   uint16
	AOE      uint8
	Damages  []AttackPair
	Position types.Vector2[int16]
}

func (a *AttackInfo) Deserialize(sr *stream.StreamReader, opcode uint16) error {
	_ = sr.Skip(1)

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
	case 2121001, 2221001, 2321001, 3221001, 3121004, 13111002, 5101004, 15101003, 5221004, 5201002:
		charge, err := sr.ReadU32()
		if err != nil {
			return err
		}
		a.Charge = charge

	default:
		a.Charge = 0
	}

	err = sr.Skip(1)
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

	if opcode == 0x1C {
		slot, err := sr.ReadU16()
		if err != nil {
			return err
		}
		a.Slot = slot

		csStar, err := sr.ReadU16()
		if err != nil {
			return err
		}
		a.CsStar = csStar

		aoe, err := sr.ReadU8()
		if err != nil {
			return err
		}
		a.AOE = aoe
	}

	damages := make([]AttackPair, 0, a.Targets)
	for range int(a.Targets) {
		oid, err := sr.ReadU32()
		if err != nil {
			return err
		}

		err = sr.Skip(14)
		if err != nil {
			return err
		}

		damagePairs := make([]DamagePair, 0, a.Hits)
		for range int(a.Hits) {
			damage, err := sr.ReadU32()
			if err != nil {
				return err
			}

			damagePairs = append(damagePairs, DamagePair{
				Damage:  damage,
				Unknown: false,
			})
		}

		damages = append(damages, AttackPair{
			OID:         oid,
			DamagePairs: damagePairs,
		})
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
		a.Position = types.Vector2[int16]{X: x, Y: y}
	}

	a.Damages = damages
	return nil
}

func (a *AttackInfo) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(0)
	sw.WriteU8((a.Targets << 4) | a.Hits)
	sw.WriteU32(a.Skill)

	switch a.Skill {
	case 2121001, 2221001, 2321001, 3221001, 3121004, 13111002, 5101004, 15101003, 5221004, 5201002:
		sw.WriteU32(a.Charge)
	default:
	}

	sw.WriteU8(0)
	sw.WriteU8(a.Unk)
	sw.WriteU8(a.Speed)
	sw.WriteU8(a.Display)
	sw.WriteU32(a.LastTick)

	if a.Slot != 0 || a.CsStar != 0 || a.AOE != 0 {
		sw.WriteU16(a.Slot)
		sw.WriteU16(a.CsStar)
		sw.WriteU8(a.AOE)
	}

	for _, damage := range a.Damages {
		sw.WriteU32(damage.OID)
		sw.Write(make([]byte, 14))

		for _, pair := range damage.DamagePairs {
			sw.WriteU32(pair.Damage)
		}
	}

	if a.Position.X != 0 || a.Position.Y != 0 {
		sw.Write16(a.Position.X)
		sw.Write16(a.Position.Y)
	}

	return nil
}
