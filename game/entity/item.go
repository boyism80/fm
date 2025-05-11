package entity

import (
	"time"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/data"
)

type Item interface {
	GetObject() *Object
	GetTemplate() data.ItemTemplate
	Serialize(writer *stream.StreamWriter, zeroPosition, leaveOut, trade bool, slot int16, itemType ItemType)
}

type BaseItem struct {
	*Object
	Template   data.ItemTemplate
	UniqueId   int64
	Expiration time.Time
}

func (item *BaseItem) GetObject() *Object {
	return item.Object
}

func (item *BaseItem) GetTemplate() data.ItemTemplate {
	return item.Template
}
