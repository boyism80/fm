package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
)

type DropMeso struct {
	Id           uint32
	Animation    constant.DropItemAnimationType
	Count        int32
	DropType     constant.DropType
	OwnerId      uint32
	SpawnedPoint types.Vector2[int16]
	IsPlayerDrop bool
	Position     types.Vector2[int16]
}

func (p *DropMeso) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0xC6)
	writer.WriteU8(uint8(p.Animation))
	writer.WriteU32(p.Id)
	writer.WriteBoolean(true)
	writer.Write32(p.Count)
	writer.WriteU32(p.OwnerId)
	writer.WriteU8(uint8(p.DropType))
	writer.Write16(p.Position.X)
	writer.Write16(p.Position.Y)
	writer.WriteU32(0)
	if p.Animation != constant.DropItemAnimationTypeNone {
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

func (p *DropMeso) Deserialize(reader *stream.StreamReader) error {
	return nil
}
