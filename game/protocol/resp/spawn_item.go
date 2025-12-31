package resp

import (
	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/core/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type SpawnItem struct {
	ID           uint32
	Animation    constant.DropItemAnimationType
	DropType     constant.DropType
	Item         entity.Item
	OwnerID      uint32
	SpawnedPoint types.Vector2[int16]
	IsPlayerDrop bool
}

func (p *SpawnItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Animation))
	writer.WriteU32(p.ID)
	writer.WriteBoolean(false)
	template := p.Item.GetModel()
	writer.WriteU32(template.GetID())
	writer.WriteU32(p.OwnerID)
	writer.WriteU8(uint8(p.DropType))
	position := p.Item.GetObject().Position
	writer.Write16(position.X)
	writer.Write16(position.Y)
	writer.WriteU32(0)
	if p.Animation != constant.DROP_ITEM_ANIMATION_TYPE_NONE {
		writer.Write16(p.SpawnedPoint.X)
		writer.Write16(p.SpawnedPoint.Y)
		writer.WriteU16(0)
	}

	writer.WriteDateTime(p.Item.GetExpiration())
	if p.IsPlayerDrop {
		writer.WriteU16(0)
	} else {
		writer.WriteU16(1)
	}

	return nil
}

func (p *SpawnItem) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *SpawnItem) Opcode() uint16 {
	return 0xC6
}
