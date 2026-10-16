package game

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

type FieldPlaceable interface {
	GetFieldPlacement() *FieldPlacement
	BindFieldPlacement(placement *FieldPlacement)
	GetCount32() int32
}

type Item interface {
	luax.Luable
	GetFieldPlacement() *FieldPlacement
	GetModel() wz.Item
	GetInventoryType() constant.InventoryType
	GetExpiration() time.Time
	SetExpiration(expiration time.Time)
	GetCount() uint16
	GetCount32() int32
	SetCount(count uint16)
	GetFlags() constant.ItemFlag
	SetFlags(flags constant.ItemFlag)
	Tradable() bool
	Storable() bool
	prepareFieldDrop() bool
	consumeKarma()
	Increase(count uint16) uint16
	Reduce(count uint16) uint16
	Clone(count uint16) Item
	BindFieldPlacement(placement *FieldPlacement)
	BuildS2C() dto.Item
	BuildRPC(ownerID uint32, slot int32) *internal.Inventory
}

type FieldPlacement struct {
	*ObjectCore
	owner                uint32
	spawnedPoint         types.Point[int16]
	dropType             constant.DropType
	playerDrop           bool
	quest                uint32
	restrictedPlayerDrop bool
	nextFFA              time.Time
	nextExpiry           time.Time
}

func (fp *FieldPlacement) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeItem
}

type fieldPlacement = FieldPlacement

type ItemCore struct {
	*fieldPlacement
	wz         wz.Item
	expiration time.Time
	count      uint16
	flags      uint16
}

type itemCore = ItemCore

type EquipmentCore struct {
	*itemCore
	enhanceChance uint8
	enhanceCount  uint8
	ownerName     string
	skillBonus    uint16
	uniqueId      *uint64
	bonusStats    *EquipmentBonusStats
}

func (e *EquipmentCore) GetEquipmentCore() *EquipmentCore { return e }
func (e *EquipmentCore) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeEquipment
}
func (e *EquipmentCore) GetCount() uint16             { return 1 }
func (e *EquipmentCore) Reduce(count uint16) uint16   { return 0 }
func (e *EquipmentCore) Increase(count uint16) uint16 { return 0 }

type Equipment interface {
	Item
	GetEquipmentCore() *EquipmentCore
	BuildEquipmentS2C() *dto.Equipment
}

func copyUint64Ptr(p *uint64) *uint64 {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

func (e *EquipmentCore) Clone(count uint16) Item {
	var bonus *EquipmentBonusStats
	if e.bonusStats != nil {
		v := *e.bonusStats
		bonus = &v
	}
	return &EquipmentCore{
		itemCore: &ItemCore{
			fieldPlacement: nil,
			count:          count,
			wz:             e.wz,
			expiration:     e.expiration,
			flags:          e.flags,
		},
		enhanceChance: e.enhanceChance,
		enhanceCount:  e.enhanceCount,
		ownerName:     e.ownerName,
		skillBonus:    e.skillBonus,
		uniqueId:      copyUint64Ptr(e.uniqueId),
		bonusStats:    bonus,
	}
}

func (item *ItemCore) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeItem
}

func (item *ItemCore) SendSpawnSyncToViewer(viewer *Character) {
	if item == nil || viewer == nil {
		return
	}
	fp := item.GetFieldPlacement()
	if fp == nil {
		return
	}
	if fp.quest > 0 && !viewer.quests.NeedsItem(fp.quest, item.GetModel().GetID()) {
		return
	}
	viewer.Send(&response.SpawnItem{
		ID:           fp.oid,
		Animation:    constant.DropItemAnimationTypeNone,
		DropType:     fp.dropType,
		ItemModel:    item.GetModel(),
		Expiration:   item.GetExpiration(),
		Position:     fp.position,
		OwnerID:      fp.owner,
		SpawnedPoint: fp.spawnedPoint,
		IsPlayerDrop: fp.playerDrop,
	}, types.SEND_POLICY_ENCRYPT)
}

func (item *ItemCore) SendDestroySyncToViewer(viewer *Character) {
	if item == nil || viewer == nil {
		return
	}
	fp := item.GetFieldPlacement()
	if fp == nil {
		return
	}
	viewer.Send(&response.RemoveItem{
		Mode:        constant.RemoveItemTypeNoAnimated,
		OID:         fp.oid,
		CharacterId: 0,
	}, types.SEND_POLICY_ENCRYPT)
}

func (item *ItemCore) GetFieldPlacement() *FieldPlacement           { return item.fieldPlacement }
func (item *ItemCore) Getcount() uint16                             { return item.count }
func (item *ItemCore) GetCount32() int32                            { return int32(item.count) }
func (item *ItemCore) SetCount(count uint16)                        { item.count = count }
func (item *ItemCore) GetFlags() constant.ItemFlag                  { return constant.ItemFlag(item.flags) }
func (item *ItemCore) SetFlags(flags constant.ItemFlag)             { item.flags = uint16(flags) }
func (item *ItemCore) GetModel() wz.Item                            { return item.wz }
func (item *ItemCore) GetExpiration() time.Time                     { return item.expiration }
func (item *ItemCore) SetExpiration(expiration time.Time)           { item.expiration = expiration }
func (item *ItemCore) BindFieldPlacement(placement *FieldPlacement) { item.fieldPlacement = placement }

func (item *ItemCore) Tradable() bool {
	if item.wz.GetTraits()&(constant.ItemTraitQuest|constant.ItemTraitPet) != 0 || item.GetFlags()&constant.ItemFlagLock != 0 {
		return false
	}
	return !item.untradeable(constant.ItemTraitTradeBlock|constant.ItemTraitAccountSharable) || item.GetFlags()&item.karma() != 0
}

func (item *ItemCore) Storable() bool {
	if item.wz.GetTraits()&constant.ItemTraitPet != 0 {
		return false
	}
	return !item.untradeable(constant.ItemTraitTradeBlock|constant.ItemTraitQuest) || item.GetFlags()&item.karma() != 0
}

func (item *ItemCore) prepareFieldDrop() bool {
	if !item.untradeable(constant.ItemTraitTradeBlock | constant.ItemTraitQuest | constant.ItemTraitAccountSharable) {
		return true
	}
	if item.GetFlags()&item.karma() == 0 {
		return false
	}
	item.consumeKarma()
	return true
}

func (item *ItemCore) untradeable(traits constant.ItemTrait) bool {
	if item.wz.GetTraits()&traits != 0 {
		return true
	}
	return item.wz.GetInventoryType() == constant.InventoryTypeEquipment && item.GetFlags()&constant.ItemFlagUntradeable != 0
}

// The client reads a different karma bit for equipment than for other tabs.
func (item *ItemCore) karma() constant.ItemFlag {
	if item.wz.GetInventoryType() == constant.InventoryTypeEquipment {
		return constant.ItemFlagEquipKarma
	}
	return constant.ItemFlagBundleKarma
}

// Karma lets an untradeable item move once; it is used up when the item arrives.
func (item *ItemCore) consumeKarma() {
	item.SetFlags(item.GetFlags() &^ item.karma())
}

func (e *EquipmentCore) bindToWearer() bool {
	if e.wz.GetTraits()&constant.ItemTraitEquipTradeBlock == 0 || e.GetFlags()&constant.ItemFlagUntradeable != 0 {
		return false
	}
	e.SetFlags(e.GetFlags() | constant.ItemFlagUntradeable)
	return true
}

func (fp *FieldPlacement) RegisterExpire(duration time.Duration) {
	fp.nextExpiry = clock.Now().Add(duration)
}
func (fp *FieldPlacement) RegisterFFA(duration time.Duration) {
	fp.nextFFA = clock.Now().Add(duration)
}
func (fp *FieldPlacement) ShouldExpire(now time.Time) bool {
	return !fp.nextExpiry.IsZero() && now.After(fp.nextExpiry)
}
func (fp *FieldPlacement) ReleaseToFFA() {
	fp.dropType = constant.DropTypeFFA
	fp.owner = 0
}

func (fp *FieldPlacement) ShouldFFA(now time.Time) bool {
	return fp.dropType != constant.DropTypeFFA && !fp.nextFFA.IsZero() && now.After(fp.nextFFA)
}

func NewItem(itemId uint32, count uint16, gw ItemWorld) (Item, error) {
	model, ok := gw.GetResources().Items[itemId]
	if !ok {
		return nil, fmt.Errorf("item model not found for ID: %d", itemId)
	}
	switch m := model.(type) {
	case wz.Equipment:
		if !m.GetEquipmentType().Equippable() {
			return nil, fmt.Errorf("item %d: model is equipment but not equippable", itemId)
		}
		return &EquipmentCore{
			itemCore: &ItemCore{
				wz:         m,
				count:      1,
				expiration: util.TimeMax,
			},
			enhanceChance: m.GetEnhanceChance(),
		}, nil
	case *wz.Consume:
		return &Consume{
			itemCore: &ItemCore{wz: m, count: count, expiration: util.TimeMax},
		}, nil
	case *wz.Installation:
		return &Installation{
			itemCore: &ItemCore{wz: m, count: 1, expiration: util.TimeMax},
		}, nil
	case *wz.MiscItem:
		return &MiscItem{
			itemCore: &ItemCore{wz: m, count: count, expiration: util.TimeMax},
		}, nil
	case *wz.CashItem:
		return &CashItem{
			itemCore: &ItemCore{wz: m, count: count, expiration: util.TimeMax},
		}, nil
	case *wz.Pet:
		uniqueID := gw.NewUniqueID()
		return &Pet{
			itemCore:    &ItemCore{wz: m, count: 1, expiration: time.Now().AddDate(0, 0, m.Life)},
			uniqueId:    &uniqueID,
			name:        gw.GetResources().GetItemName(itemId),
			level:       1,
			fullness:    constant.PetMaxFullness,
			skills:      m.Skills,
			secondsLeft: uint32(m.LimitedLife),
		}, nil
	default:
		return nil, fmt.Errorf("not supported item type")
	}
}
