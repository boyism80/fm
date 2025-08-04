package resp

import (
	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/core/types"
	"github.com/boyism80/fm/game/action"
	"github.com/boyism80/fm/game/entity"
)

type Move struct {
	Character  *entity.Character
	Fragments  []action.MoveFragment
	StartPoint types.Vector2[int16]
}

func (p *Move) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Character.ID)
	writer.Write16(p.Character.Position.X)
	writer.Write16(p.Character.Position.Y)
	writer.WriteU8(uint8(len(p.Fragments)))
	for _, move := range p.Fragments {
		move.Serialize(writer)
	}
	return nil
}

func (p *Move) Opcode() uint16 {
	return 0x82 // Move opcode
}

func (p *Move) Deserialize(reader *stream.StreamReader) error {
	return nil
}
