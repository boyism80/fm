package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/game/constant"
)

type SpawnMeso struct {
	ID           uint32
	Animation    constant.DropItemAnimationType
	Count        int32
	DropType     constant.DropType
	OwnerID      uint32
	SpawnedPoint types.Vector2[int16]
	IsPlayerDrop bool
	Position     types.Vector2[int16]
}

func (p *SpawnMeso) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Animation))
	writer.WriteU32(p.ID)
	writer.WriteBoolean(true)
	writer.Write32(p.Count)
	writer.WriteU32(p.OwnerID)
	writer.WriteU8(uint8(p.DropType))
	writer.Write16(p.Position.X)
	writer.Write16(p.Position.Y)
	writer.WriteU32(0)
	if p.Animation != constant.DROP_ITEM_ANIMATION_TYPE_NONE {
		writer.Write16(p.SpawnedPoint.X)
		writer.Write16(p.SpawnedPoint.Y)
		writer.WriteU16(0)
	}

	if p.IsPlayerDrop {
		writer.WriteU16(0)
	} else {
		writer.WriteU16(1)
	}

	return nil
}

func (p *SpawnMeso) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *SpawnMeso) Opcode() uint16 {
	return 0xC6
}
