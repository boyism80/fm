package response

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type SpawnItem struct {
	ID           uint32
	Animation    constant.DropItemAnimationType
	DropType     constant.DropType
	ItemModel    wz.Item
	Expiration   time.Time
	Position     types.Vector2[int16]
	OwnerID      uint32
	SpawnedPoint types.Vector2[int16]
	IsPlayerDrop bool
}

func (p *SpawnItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Animation))
	writer.WriteU32(p.ID)
	writer.WriteBoolean(false)
	if p.ItemModel != nil {
		writer.WriteU32(p.ItemModel.GetID())
	} else {
		writer.WriteU32(0)
	}
	writer.WriteU32(p.OwnerID)
	writer.WriteU8(uint8(p.DropType))
	writer.Write16(p.Position.X)
	writer.Write16(p.Position.Y)
	writer.WriteU32(0)
	if p.Animation != constant.DropItemAnimationTypeNone {
		writer.Write16(p.SpawnedPoint.X)
		writer.Write16(p.SpawnedPoint.Y)
		writer.WriteU16(0)
	}

	if !p.Expiration.IsZero() {
		writer.WriteDateTime(p.Expiration)
	} else {
		writer.WriteU64(0)
	}
	if p.IsPlayerDrop {
		writer.WriteU16(0)
	} else {
		writer.WriteU16(1)
	}

	return nil
}

func (p *SpawnItem) Deserialize(reader *stream.StreamReader) {
	p.Animation = constant.DropItemAnimationType(reader.ReadU8())
	p.ID = reader.ReadU32()
	reader.Skip(1)
	p.ItemModel = &wz.ItemCore{ID: reader.ReadU32()}
	p.OwnerID = reader.ReadU32()
	p.DropType = constant.DropType(reader.ReadU8())
	p.Position = types.Vector2[int16]{X: reader.Read16(), Y: reader.Read16()}
	reader.Skip(4)
	if p.Animation != constant.DropItemAnimationTypeNone {
		p.SpawnedPoint = types.Vector2[int16]{X: reader.Read16(), Y: reader.Read16()}
		reader.Skip(2)
	}
	reader.Skip(8)
	p.IsPlayerDrop = reader.ReadU16() == 0
}

func (p *SpawnItem) Opcode() uint16 {
	return 0xC6
}
