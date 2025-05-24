package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/common/context"
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/data"
)

type Item interface {
	GetObject() *Object
	GetSpec() data.ItemSpec
	GetInventoryType() constant.InventoryType
	GetExpiration() time.Time
	GetCount() uint16
	Serialize(writer *stream.StreamWriter, trade bool, slot int16)
}

type ItemCore struct {
	*Object
	Spec       data.ItemSpec
	UniqueId   int64
	Expiration time.Time
}

type CashItem struct {
	*ItemCore
	Count     uint16
	OwnerName string
	Flags     uint16
}

type Installation struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

type GeneralItem struct {
	*ItemCore
	Count     uint16
	OwnerName string
	Flags     uint16
}

type Equipment struct {
	*ItemCore
	EnchantChance uint8
	OwnerName     string
	Flag          uint16
	SkillBonus    uint16
}

type Consume struct {
	*ItemCore
	Count     uint16
	OwnerName string
	Flags     uint16
}

type Pet struct {
	*ItemCore
	Level       uint8
	Closeness   uint16
	Fullness    uint8
	Speed       uint16
	Flags       uint16
	SecondsLeft uint32
	Expiration  time.Time
}

func (item *ItemCore) GetObject() *Object {
	return item.Object
}

func (item *ItemCore) GetSpec() data.ItemSpec {
	return item.Spec
}

func (item *ItemCore) GetExpiration() time.Time {
	return item.Expiration
}

func (item *CashItem) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeCash
}

func (item *CashItem) GetCount() uint16 {
	return item.Count
}

func (item *CashItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(constant.ItemTypeEtc))
	writer.WriteU32(item.Spec.GetID())

	hasUID := item.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(item.UniqueId)
	}

	writer.WriteDateTime(item.Expiration)
	writer.WriteU16(item.Count)
	writer.WriteStr16(item.OwnerName)
	writer.WriteU16(item.Flags)
}

func (item *Installation) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeInstallation
}

func (item *Installation) GetCount() uint16 {
	return 1
}

func (item *Installation) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(constant.ItemTypeEtc))
	writer.WriteU32(item.Spec.GetID())

	hasUID := item.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(item.UniqueId)
	}

	writer.WriteDateTime(item.Expiration)
	writer.WriteU16(1)
	writer.WriteStr16(item.OwnerName)
	writer.WriteU16(item.Flags)
}

func (item *GeneralItem) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeEtc
}

func (item *GeneralItem) GetCount() uint16 {
	return item.Count
}

func (item *GeneralItem) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(constant.ItemTypeEtc))
	writer.WriteU32(item.Spec.GetID())

	hasUID := item.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(item.UniqueId)
	}

	writer.WriteDateTime(item.Expiration)
	writer.WriteU16(item.Count)
	writer.WriteStr16(item.OwnerName)
	writer.WriteU16(item.Flags)
}

func (item *Equipment) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeEquipment
}

func (item *Equipment) GetCount() uint16 {
	return 1
}

func (equipment *Equipment) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	spec, _ := equipment.Spec.(*data.EquipmentSpec)
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

	writer.WriteU8(uint8(constant.ItemTypeEquipment))
	writer.WriteU32(spec.Id)

	hasUID := equipment.UniqueId > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(equipment.UniqueId)
	}

	writer.WriteDateTime(equipment.Expiration)
	writer.WriteU8(equipment.EnchantChance)
	writer.WriteU8(spec.Required.Level)
	writer.WriteU16(spec.Ability.Str)
	writer.WriteU16(spec.Ability.Dex)
	writer.WriteU16(spec.Ability.Int)
	writer.WriteU16(spec.Ability.Luk)
	writer.WriteU16(spec.Ability.MaxHP)
	writer.WriteU16(spec.Ability.MaxMP)
	writer.WriteU16(spec.Ability.PAD)
	writer.WriteU16(spec.Ability.MAD)
	writer.WriteU16(spec.Ability.PDD)
	writer.WriteU16(spec.Ability.MDD)
	writer.WriteU16(spec.Ability.ACC)
	writer.WriteU16(spec.Ability.Avoid)
	writer.WriteU16(spec.Ability.Hands)
	writer.WriteU16(spec.Ability.Speed)
	writer.WriteU16(spec.Ability.Jump)
	writer.WriteStr16(equipment.OwnerName)
	writer.WriteU16(equipment.Flag)
	writer.WriteBoolean(equipment.SkillBonus > 0)
	writer.WriteU8(1)
	writer.WriteU32(0)
	if equipment.UniqueId <= 0 {
		inventoryId := 0
		if inventoryId > 0 {
			writer.WriteU64(uint64(inventoryId))
		} else {
			writer.Write64(-1)
		}
	}
	writer.WriteDateTime(util.TimeZero)
	writer.Write32(-1)
}

func (item *Consume) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeConsume
}

func (item *Consume) GetCount() uint16 {
	return item.Count
}

func (item *Consume) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	writer.WriteU8(uint8(slot))
	writer.WriteU8(uint8(constant.ItemTypeEtc))
	writer.WriteU32(item.Spec.GetID())

	hasUID := item.UniqueId > 0
	writer.WriteBoolean(false)
	if hasUID {
		writer.Write64(item.UniqueId)
	}

	writer.WriteDateTime(item.Expiration)
	spec, ok := item.Spec.(*data.ConsumeSpec)
	if !ok {
		return
	}

	writer.WriteU16(item.Count)
	writer.WriteStr16(item.OwnerName)
	writer.WriteU16(item.Flags)

	isThrowingStart := spec.Id/10000 == 207
	isBullet := spec.Id/10000 == 233
	isWhat := spec.Id/10000 == 287
	inventoryId := uint64(54399043)
	if isThrowingStart || isBullet || isWhat {
		writer.WriteU64(inventoryId)
	}
}

func (item *Pet) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeCash
}

func (item *Pet) GetCount() uint16 {
	return 1
}

func (pet *Pet) Serialize(writer *stream.StreamWriter, trade bool, slot int16) {
	spec, ok := pet.Spec.(*data.PetSpec)
	if !ok {
		return
	}

	writer.WriteU8(uint8(slot))
	writer.WriteU8(3)
	writer.WriteU32(spec.Id)

	hasUID := pet.UniqueId > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(pet.UniqueId)
	}

	writer.WriteDateTime(pet.ItemCore.Expiration)
	writer.WriteStaticStr(spec.Name, 13)
	writer.WriteU8(pet.Level)
	writer.WriteU16(pet.Closeness)
	writer.WriteU8(pet.Fullness)
	writer.WriteDateTime(pet.Expiration)
	writer.WriteU16(pet.Speed)
	writer.WriteU16(pet.Flags)
	if spec.Id == 5000054 && pet.SecondsLeft > 0 {
		writer.WriteU32(pet.SecondsLeft)
	} else {
		writer.WriteU32(0)
	}
}

func NewItem(ctx *context.ServerContext, itemId uint32, count uint16) (Item, error) {
	t, ok := ctx.Resources.Items[uint32(itemId)]
	if !ok {
		return nil, fmt.Errorf("%d is not valid item id", itemId)
	}

	switch spec := t.(type) {
	case *data.EquipmentSpec:
		return &Equipment{
			ItemCore: &ItemCore{
				Spec: spec,
			},
			EnchantChance: spec.TUC,
		}, nil

	case *data.ConsumeSpec:
		return &Consume{
			ItemCore: &ItemCore{
				Spec: spec,
			},
			Count: count,
		}, nil

	case *data.InstallationSpec:
		return &Installation{
			ItemCore: &ItemCore{
				Spec: spec,
			},
		}, nil

	case *data.GeneralItemSpec:
		return &GeneralItem{
			ItemCore: &ItemCore{
				Spec: spec,
			},
			Count: count,
		}, nil

	case *data.CashItemSpec:
		return &CashItem{
			ItemCore: &ItemCore{
				Spec: spec,
			},
			Count: count,
		}, nil

	case *data.PetSpec:
		petExpiration, err := time.ParseInLocation("2006-01-02 15:04:05", "2025-05-30 09:30:00", util.KST)
		if err != nil {
			fmt.Println(err)
		}
		return &Pet{
			ItemCore: &ItemCore{
				Spec: spec,
			},
			Expiration: petExpiration,
		}, nil

	default:
		return nil, fmt.Errorf("not supported item type")
	}
}
