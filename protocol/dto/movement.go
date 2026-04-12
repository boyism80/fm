package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

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

func ReadMovements(reader *stream.StreamReader) []MoveFragment {
	numCommands := reader.ReadU8()

	fragments := make([]MoveFragment, 0, numCommands)

	for range int(numCommands) {
		cmd := reader.ReadU8()

		switch cmd {
		case 0, 5, 17:
			x := reader.Read16()
			y := reader.Read16()
			vx := reader.Read16()
			vy := reader.Read16()
			foothold := reader.Read16()
			stance := reader.ReadU8()
			duration := reader.Read16()

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
			x := reader.Read16()
			y := reader.Read16()
			vx := reader.Read16()
			vy := reader.Read16()
			unk := reader.Read16()
			foothold := reader.Read16()
			stance := reader.ReadU8()
			duration := reader.Read16()
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
			x := reader.Read16()
			y := reader.Read16()
			stance := reader.ReadU8()
			duration := reader.Read16()

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
			x := reader.Read16()
			y := reader.Read16()
			vx := reader.Read16()
			vy := reader.Read16()
			stance := reader.ReadU8()

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
			stance := reader.ReadU8()
			frag := NoneMovement{
				BasicMovement: &BasicMovement{
					Command: cmd,
					Stance:  stance,
				},
			}
			fragments = append(fragments, &frag)

		case 14:
			x := reader.Read16()
			y := reader.Read16()
			foothold := reader.Read16()
			stance := reader.ReadU8()
			duration := reader.Read16()

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
			stance := reader.ReadU8()
			foothold := reader.Read16()
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

	return fragments
}
