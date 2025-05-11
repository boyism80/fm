package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/data"
)

type Item interface {
	GetObject() *Object
	GetTemplate() data.ItemTemplate
	GetInventoryType() InventoryType
	Serialize(writer *stream.StreamWriter, trade bool, slot int16)
}

func NewItem(ctx *context.ServerContext, itemId uint32, count uint16) (Item, error) {
	t, ok := ctx.GameData.Items[uint32(itemId)]
	if !ok {
		return nil, fmt.Errorf("%d is not valid item id", itemId)
	}

	switch template := t.(type) {
	case *data.EquipmentTemplate:
		return &Equipment{
			BaseItem: &BaseItem{
				Template: template,
			},
			EnchantChance: template.TUC,
		}, nil

	case *data.ConsumeTemplate:
		return &Consume{
			BaseItem: &BaseItem{
				Template: template,
			},
			Count: count,
		}, nil

	case *data.InstallationTemplate:
		return &Installation{
			BaseItem: &BaseItem{
				Template: template,
			},
		}, nil

	case *data.GeneralItemTemplate:
		return &GeneralItem{
			BaseItem: &BaseItem{
				Template: template,
			},
			Count: count,
		}, nil

	case *data.CashItemTemplate:
		return &CashItem{
			BaseItem: &BaseItem{
				Template: template,
			},
			Count: count,
		}, nil

	case *data.PetTemplate:
		petExpiration, err := time.ParseInLocation("2006-01-02 15:04:05", "2025-05-30 09:30:00", util.KST)
		if err != nil {
			fmt.Println(err)
		}
		return &Pet{
			BaseItem: &BaseItem{
				Template: template,
			},
			Expiration: petExpiration,
		}, nil

	default:
		return nil, fmt.Errorf("not supported item type")
	}
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
