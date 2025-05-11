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

type baseItem struct {
	*Object
	Template   data.ItemTemplate
	UniqueId   int64
	Expiration time.Time
}

func (item *baseItem) GetObject() *Object {
	return item.Object
}

func (item *baseItem) GetTemplate() data.ItemTemplate {
	return item.Template
}

func NewItem(ctx *context.ServerContext, itemId uint32, count uint16) (Item, error) {
	t, ok := ctx.Resources.Items[uint32(itemId)]
	if !ok {
		return nil, fmt.Errorf("%d is not valid item id", itemId)
	}

	switch template := t.(type) {
	case *data.EquipmentTemplate:
		return &Equipment{
			baseItem: &baseItem{
				Template: template,
			},
			EnchantChance: template.TUC,
		}, nil

	case *data.ConsumeTemplate:
		return &Consume{
			baseItem: &baseItem{
				Template: template,
			},
			Count: count,
		}, nil

	case *data.InstallationTemplate:
		return &Installation{
			baseItem: &baseItem{
				Template: template,
			},
		}, nil

	case *data.GeneralItemTemplate:
		return &GeneralItem{
			baseItem: &baseItem{
				Template: template,
			},
			Count: count,
		}, nil

	case *data.CashItemTemplate:
		return &CashItem{
			baseItem: &baseItem{
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
			baseItem: &baseItem{
				Template: template,
			},
			Expiration: petExpiration,
		}, nil

	default:
		return nil, fmt.Errorf("not supported item type")
	}
}

type Equipment struct {
	*baseItem
	EnchantChance uint8
	OwnerName     string
	Flag          uint16
	SkillBonus    uint16
}

func (equipment *Equipment) GetInventoryType() InventoryType {
	return InventoryTypeEquip
}

func (equipment *Equipment) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	template, ok := equipment.Template.(*data.EquipmentTemplate)
	if !ok {
		return // TODO: return error
	}

	if slot <= -1 {
		slot *= -1
		if slot > 100 && slot < 1000 {
			slot -= 100
		}
	}
	if slot != 0 && !trade {
		writer.WriteU16(uint16(slot))
	} else {
		writer.WriteU8(uint8(slot))
	}

	writer.WriteU8(uint8(ItemTypeEquipment))
	writer.WriteU32(template.Id)

	hasUID := equipment.UniqueId > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(equipment.UniqueId)
	}

	writer.WriteDateTime(equipment.Expiration)
	writer.WriteU8(equipment.EnchantChance)
	writer.WriteU8(template.Required.Level)
	writer.WriteU16(template.Ability.Str)
	writer.WriteU16(template.Ability.Dex)
	writer.WriteU16(template.Ability.Int)
	writer.WriteU16(template.Ability.Luk)
	writer.WriteU16(template.Ability.MaxHP)
	writer.WriteU16(template.Ability.MaxMP)
	writer.WriteU16(template.Ability.PAD)
	writer.WriteU16(template.Ability.MAD)
	writer.WriteU16(template.Ability.PDD)
	writer.WriteU16(template.Ability.MDD)
	writer.WriteU16(template.Ability.ACC)
	writer.WriteU16(template.Ability.Avoid)
	writer.WriteU16(template.Ability.Hands)
	writer.WriteU16(template.Ability.Speed)
	writer.WriteU16(template.Ability.Jump)
	writer.WriteStr16(equipment.OwnerName)
	writer.WriteU16(equipment.Flag)
	writer.WriteBoolean(equipment.SkillBonus > 0)
	writer.WriteU8(1)  // item level?
	writer.WriteU32(0) // item exp percent?
	if equipment.UniqueId <= 0 {
		inventoryId := 0 // TODO: tracking
		if inventoryId > 0 {
			writer.WriteU64(uint64(inventoryId))
		} else {
			writer.Write64(-1)
		}
	}
	writer.WriteDateTime(util.TimeZero)
	writer.Write32(-1)
}
