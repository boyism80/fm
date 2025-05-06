package entity

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/data"
)

type Item interface {
	GetTemplate() data.ItemTemplate
	GetParts() uint32
	Serialize(writer *stream.StreamWriter, zeroPosition, leaveOut, trade, bagSlot bool)
}

type BaseItem struct {
	Template data.ItemTemplate
	UniqueId int64
}

func (item *BaseItem) GetTemplate() data.ItemTemplate {
	return item.Template
}
