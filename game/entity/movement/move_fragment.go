package movement

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
)

type MoveFragment interface {
	GetNewState() uint8
	Serialize(sw *stream.StreamWriter)
}

func Parse(reader *stream.StreamReader) ([]MoveFragment, error) {
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
				Command:  cmd,
				Position: types.Vec2[int16]{X: x, Y: y},
				Foothold: foothold,
				Stance:   stance,
				Duration: duration,
				Velocity: types.Vec2[int16]{X: vx, Y: vy},
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
				Command:  cmd,
				Position: types.Vec2[int16]{X: x, Y: y},
				Duration: duration,
				Stance:   stance,
				Velocity: types.Vec2[int16]{X: vx, Y: vy},
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
				Command:  cmd,
				Position: types.Vec2[int16]{X: x, Y: y},
				Stance:   stance,
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
				Command:  cmd,
				Position: types.Vec2[int16]{X: x, Y: y},
				Velocity: types.Vec2[int16]{X: vx, Y: vy},
				Stance:   stance,
			}

			fragments = append(fragments, &frag)

		case 10:
			stance, _ := reader.ReadU8()
			frag := NoneMovement{
				Command: cmd,
				Stance:  stance,
			}
			fragments = append(fragments, &frag)

		case 14:
			x, _ := reader.Read16()
			y, _ := reader.Read16()
			foothold, _ := reader.Read16()
			stance, _ := reader.ReadU8()
			duration, _ := reader.Read16()

			frag := ChairMovement{
				Command:  cmd,
				Position: types.Vec2[int16]{X: x, Y: y},
				Foothold: foothold,
				Stance:   stance,
				Duration: duration,
			}
			fragments = append(fragments, &frag)

		default:
			stance, _ := reader.ReadU8()
			foothold, _ := reader.Read16()
			frag := AranMovement{
				Command:  cmd,
				Stance:   stance,
				Foothold: foothold,
				Position: types.Vec2[int16]{},
			}
			fragments = append(fragments, &frag)
		}
	}

	return fragments, nil
}
