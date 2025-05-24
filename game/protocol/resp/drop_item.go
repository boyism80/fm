package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type DropItem struct {
	Id           uint32
	Animation    constant.DropItemAnimationType
	Meso         uint32
	DropType     uint8
	Item         entity.Item
	OwnerId      uint32
	DropFrom     types.Vec2[int16]
	IsPlayerDrop bool
}

func (p *DropItem) Serialize(writer *stream.StreamWriter) error {
	isMeso := p.Meso > 0
	writer.WriteU16(0xC6)
	writer.WriteU8(uint8(p.Animation))
	writer.WriteU32(p.Id)
	writer.WriteBoolean(isMeso)
	if isMeso {
		writer.WriteU32(p.Meso)
	} else {
		template := p.Item.GetSpec()
		writer.WriteU32(template.GetID())
	}
	writer.WriteU32(p.OwnerId)
	writer.WriteU8(p.DropType)
	position := p.Item.GetObject().Position
	writer.Write16(position.X)
	writer.Write16(position.Y)
	writer.WriteU32(0)
	if p.Animation != constant.DropItemAnimationTypeNone {
		writer.Write16(p.DropFrom.X)
		writer.Write16(p.DropFrom.Y)
		writer.WriteU16(0)
	}

	if !isMeso {
		writer.WriteDateTime(p.Item.GetExpiration())
	}

	if p.IsPlayerDrop {
		writer.WriteU16(0)
	} else {
		writer.WriteU16(1)
	}

	return nil
}

func (p *DropItem) Deserialize(reader *stream.StreamReader) error {
	return nil
}
