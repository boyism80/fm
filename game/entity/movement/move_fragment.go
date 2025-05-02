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
			xpos, _ := reader.Read16()
			ypos, _ := reader.Read16()
			xwobble, _ := reader.Read16()
			ywobble, _ := reader.Read16()
			unk, _ := reader.Read16()
			newstate, _ := reader.ReadU8()
			duration, _ := reader.Read16()

			frag := AbsoluteLifeMovement{
				Command:        cmd,
				Position:       types.Vec2[int16]{X: xpos, Y: ypos},
				Unknown:        unk,
				NewState:       newstate,
				Duration:       duration,
				PixelPerSecond: types.Vec2[int16]{X: xwobble, Y: ywobble},
			}

			fragments = append(fragments, &frag)

		case 15:
			xpos, _ := reader.Read16()
			ypos, _ := reader.Read16()
			xwobble, _ := reader.Read16()
			ywobble, _ := reader.Read16()
			unk, _ := reader.Read16()
			fh, _ := reader.Read16()
			newstate, _ := reader.ReadU8()
			duration, _ := reader.Read16()
			frag := JumpDownMovement{
				Command:        cmd,
				Position:       types.Vec2[int16]{X: xpos, Y: ypos},
				Duration:       duration,
				NewState:       newstate,
				PixelPerSecond: types.Vec2[int16]{X: xwobble, Y: ywobble},
				Unknown:        unk,
				FH:             fh,
			}

			fragments = append(fragments, &frag)

		case 1, 2, 6, 12, 13, 16:
			xmod, _ := reader.Read16()
			ymod, _ := reader.Read16()
			newstate, _ := reader.ReadU8()
			duration, _ := reader.Read16()

			frag := RelativeLifeMovement{
				Command:  cmd,
				Position: types.Vec2[int16]{X: xmod, Y: ymod},
				NewState: newstate,
				Duration: duration,
			}

			fragments = append(fragments, &frag)

		case 3, 4, 7, 8, 9, 11:
			xpos, _ := reader.Read16()
			ypos, _ := reader.Read16()
			xwobble, _ := reader.Read16()
			ywobble, _ := reader.Read16()
			newstate, _ := reader.ReadU8()

			frag := TeleportMovement{
				Command:        cmd,
				Position:       types.Vec2[int16]{X: xpos, Y: ypos},
				PixelPerSecond: types.Vec2[int16]{X: xwobble, Y: ywobble},
				NewState:       newstate,
			}

			fragments = append(fragments, &frag)

		case 10:
			news, _ := reader.ReadU8()
			frag := NoneMovement{
				Command:  cmd,
				NewState: news,
			}
			fragments = append(fragments, &frag)

		case 14:
			xpos, _ := reader.Read16()
			ypos, _ := reader.Read16()
			unk, _ := reader.Read16()
			newstate, _ := reader.ReadU8()
			duration, _ := reader.Read16()

			frag := ChairMovement{
				Command:  cmd,
				Position: types.Vec2[int16]{X: xpos, Y: ypos},
				Unknown:  unk,
				NewState: newstate,
				Duration: duration,
			}
			fragments = append(fragments, &frag)

		default:
			newstate, _ := reader.ReadU8()
			unk, _ := reader.Read16()
			frag := AranMovement{
				Command:  cmd,
				NewState: newstate,
				Unknown:  unk,
				Position: types.Vec2[int16]{},
			}
			fragments = append(fragments, &frag)
		}
	}

	return fragments, nil
}
