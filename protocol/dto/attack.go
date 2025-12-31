package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

// Attack types
type DamagePair struct {
	Damage  uint32
	Unknown bool
}

type AttackPair struct {
	OID         uint32
	Position    types.Vector2[int16]
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
	Slot     uint8 // only for RANGED_ATTACK
	CsStar   uint8 // only for RANGED_ATTACK
	AOE      uint8 // only for RANGED_ATTACK
	Damages  []AttackPair
	Position types.Vector2[int16]
}

// Deserialize deserializes AttackInfo from stream
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

	if opcode == 0x1C {
		slot, err := sr.ReadU8()
		if err != nil {
			return err
		}
		a.Slot = slot

		csStar, err := sr.ReadU8()
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

	x, err := sr.Read16()
	if err != nil {
		return err
	}
	y, err := sr.Read16()
	if err != nil {
		return err
	}
	a.Position = types.Vector2[int16]{X: x, Y: y}

	damages := make([]AttackPair, 0, a.Targets)
	for range int(a.Targets) {
		oid, err := sr.ReadU32()
		if err != nil {
			return err
		}

		x, err := sr.Read16()
		if err != nil {
			return err
		}
		y, err := sr.Read16()
		if err != nil {
			return err
		}

		damagePairs := make([]DamagePair, 0, a.Hits)
		for range int(a.Hits) {
			damage, err := sr.ReadU32()
			if err != nil {
				return err
			}

			unknown, err := sr.ReadBool()
			if err != nil {
				return err
			}

			damagePairs = append(damagePairs, DamagePair{
				Damage:  damage,
				Unknown: unknown,
			})
		}

		damages = append(damages, AttackPair{
			OID:         oid,
			Position:    types.Vector2[int16]{X: x, Y: y},
			DamagePairs: damagePairs,
		})
	}

	a.Damages = damages
	return nil
}

// Serialize serializes AttackInfo to stream
func (a *AttackInfo) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(0) // unknown byte
	sw.WriteU8((a.Targets << 4) | a.Hits)
	sw.WriteU32(a.Skill)

	switch a.Skill {
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
		sw.WriteU32(a.Charge)
	default:
		// Charge is already 0
	}

	sw.WriteU8(0) // nOption
	sw.WriteU8(a.Unk)
	sw.WriteU8(a.Speed)
	sw.WriteU8(a.Display)
	sw.WriteU32(a.LastTick)

	if a.Slot != 0 || a.CsStar != 0 || a.AOE != 0 {
		sw.WriteU8(a.Slot)
		sw.WriteU8(a.CsStar)
		sw.WriteU8(a.AOE)
	}

	sw.Write16(a.Position.X)
	sw.Write16(a.Position.Y)

	for _, damage := range a.Damages {
		sw.WriteU32(damage.OID)
		sw.Write16(damage.Position.X)
		sw.Write16(damage.Position.Y)

		for _, pair := range damage.DamagePairs {
			sw.WriteU32(pair.Damage)
			sw.WriteBoolean(pair.Unknown)
		}
	}

	return nil
}
