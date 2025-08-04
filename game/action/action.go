package action

import (
	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/core/types"
)

// Movement types
type MoveFragment interface {
	GetStance() uint8
	Serialize(sw *stream.StreamWriter)
}

type BasicMovement struct {
	Command uint8
	Stance  uint8
}

type AbsoluteLifeMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Duration int16
	Velocity types.Vector2[int16]
	Offset   types.Vector2[int16]
	Foothold int16
}

type AranMovement struct {
	*BasicMovement
	Foothold int16
	Position types.Vector2[int16]
}

type ChairMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Foothold int16
	Duration int16
}

type JumpDownMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Duration int16
	Velocity types.Vector2[int16]
	Unknown  int16
	Foothold int16
}

type NoneMovement struct {
	*BasicMovement
}

type RelativeLifeMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Duration int16
}

type TeleportMovement struct {
	*BasicMovement
	Position types.Vector2[int16]
	Velocity types.Vector2[int16]
}

// Movement methods
func (m *BasicMovement) GetStance() uint8 {
	return m.Stance
}

func (m *AbsoluteLifeMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Velocity.X)
	writer.Write16(m.Velocity.Y)
	writer.Write16(m.Foothold)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *AranMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Foothold)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Foothold)
}

func (m *ChairMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Foothold)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *JumpDownMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Velocity.X)
	writer.Write16(m.Velocity.Y)
	writer.Write16(m.Unknown)
	writer.Write16(m.Foothold)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *NoneMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.WriteU8(m.Stance)
}

func (m *RelativeLifeMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Duration)
}

func (m *TeleportMovement) Serialize(writer *stream.StreamWriter) {
	writer.WriteU8(m.Command)
	writer.Write16(m.Position.X)
	writer.Write16(m.Position.Y)
	writer.Write16(m.Velocity.X)
	writer.WriteU8(m.Stance)
	writer.Write16(m.Velocity.Y)
}

func ReadMovements(reader *stream.StreamReader) ([]MoveFragment, error) {
	numCommands, err := reader.ReadU8()
	if err != nil {
		return nil, err
	}

	fragments := make([]MoveFragment, 0, numCommands)

	for range int(numCommands) {
		cmd, err := reader.ReadU8()
		if err != nil {
			return nil, err
		}

		switch cmd {
		case 0, 5, 17:
			x, _ := reader.Read16()
			y, _ := reader.Read16()
			vx, _ := reader.Read16()
			vy, _ := reader.Read16()
			foothold, _ := reader.Read16()
			stance, _ := reader.ReadU8()
			duration, _ := reader.Read16()

			frag := AbsoluteLifeMovement{
				BasicMovement: &BasicMovement{
					Command: cmd,
					Stance:  stance,
				},
				Position: types.Vector2[int16]{X: x, Y: y},
				Foothold: foothold,
				Duration: duration,
				Velocity: types.Vector2[int16]{X: vx, Y: vy},
			}

			fragments = append(fragments, &frag)

		case 15:
			x, _ := reader.Read16()
			y, _ := reader.Read16()
			vx, _ := reader.Read16()
			vy, _ := reader.Read16()
			unk, _ := reader.Read16()
			foothold, _ := reader.Read16()
			stance, _ := reader.ReadU8()
			duration, _ := reader.Read16()
			frag := JumpDownMovement{
				BasicMovement: &BasicMovement{
					Command: cmd,
					Stance:  stance,
				},
				Position: types.Vector2[int16]{X: x, Y: y},
				Duration: duration,
				Velocity: types.Vector2[int16]{X: vx, Y: vy},
				Unknown:  unk,
				Foothold: foothold,
			}

			fragments = append(fragments, &frag)

		case 1, 2, 6, 12, 13, 16:
			x, _ := reader.Read16()
			y, _ := reader.Read16()
			stance, _ := reader.ReadU8()
			duration, _ := reader.Read16()

			frag := RelativeLifeMovement{
				BasicMovement: &BasicMovement{
					Command: cmd,
					Stance:  stance,
				},
				Position: types.Vector2[int16]{X: x, Y: y},
				Duration: duration,
			}

			fragments = append(fragments, &frag)

		case 3, 4, 7, 8, 9, 11:
			x, _ := reader.Read16()
			y, _ := reader.Read16()
			vx, _ := reader.Read16()
			vy, _ := reader.Read16()
			stance, _ := reader.ReadU8()

			frag := TeleportMovement{
				BasicMovement: &BasicMovement{
					Command: cmd,
					Stance:  stance,
				},
				Position: types.Vector2[int16]{X: x, Y: y},
				Velocity: types.Vector2[int16]{X: vx, Y: vy},
			}

			fragments = append(fragments, &frag)

		case 10:
			stance, _ := reader.ReadU8()
			frag := NoneMovement{
				BasicMovement: &BasicMovement{
					Command: cmd,
					Stance:  stance,
				},
			}
			fragments = append(fragments, &frag)

		case 14:
			x, _ := reader.Read16()
			y, _ := reader.Read16()
			foothold, _ := reader.Read16()
			stance, _ := reader.ReadU8()
			duration, _ := reader.Read16()

			frag := ChairMovement{
				BasicMovement: &BasicMovement{
					Command: cmd,
					Stance:  stance,
				},
				Position: types.Vector2[int16]{X: x, Y: y},
				Foothold: foothold,
				Duration: duration,
			}
			fragments = append(fragments, &frag)

		default:
			stance, _ := reader.ReadU8()
			foothold, _ := reader.Read16()
			frag := AranMovement{
				BasicMovement: &BasicMovement{
					Command: cmd,
					Stance:  stance,
				},
				Foothold: foothold,
				Position: types.Vector2[int16]{},
			}
			fragments = append(fragments, &frag)
		}
	}

	return fragments, nil
}

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

// Attack methods
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

	a.Damages = make([]AttackPair, 0, a.Targets)
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

		a.Damages = append(a.Damages, AttackPair{
			OID:         oid,
			DamagePairs: attack,
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
		a.Position = types.Vector2[int16]{X: x, Y: y}
	}

	return nil
}
