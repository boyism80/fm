package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/entity/movement"
)

type Move struct {
	Character     *entity.Character
	MoveFragments []movement.MoveFragment
	StartPoint    types.Vec2[int16]
}

func (a *Move) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x82)
	writer.WriteU32(a.Character.Id)
	writer.Write16(a.Character.Position.X)
	writer.Write16(a.Character.Position.Y)
	writer.WriteU8(uint8(len(a.MoveFragments)))
	for _, move := range a.MoveFragments {
		move.Serialize(writer)
	}
	return nil
}

func (a *Move) Deserialize(reader *stream.StreamReader) error {
	return nil
}
