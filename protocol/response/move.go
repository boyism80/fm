package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type Move struct {
	Character  *dto.Character
	Fragments  []dto.MoveFragment
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
	return 0x82
}

func (p *Move) Deserialize(reader *stream.StreamReader) {
}
