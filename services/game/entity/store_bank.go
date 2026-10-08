package entity

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

const (
	StoreBankNothingMapID         = 999999999
	StoreBankNothingChannel       = 0xFF
	StoreBankSlots                = 16
	StoreBankFeeGrace             = 24 * time.Hour
	StoreBankFeePercentMax  int64 = 100
)

var (
	ErrStoreBankClosed = errors.New("store bank is not open")
	ErrStoreBankBusy   = errors.New("store bank request is pending")
)

type StoreBank struct {
	owner    *Character
	shops    []*internal.Shop
	closedAt time.Time
	Meso     int32
	Items    []*ShopItem
	pending  bool
}

func (sb *StoreBank) Open(actx actor.Context, npcID uint32) {
	if sb.pending {
		return
	}

	sb.pending = true
	sb.owner.Listener.FindEntrustedShopAsync(actx, sb.owner).Do(func(v *internal.FindEntrustedShopReply) error {
		sb.pending = false
		var shops []*internal.Shop
		for _, closed := range v.GetStoreBank() {
			if closed.GetCharacterId() == sb.owner.GetID() {
				shops = append(shops, closed)
			}
		}
		shop := v.GetShop()
		switch {
		case len(shops) > 0:
			sb.load(shops)
			sb.owner.Listener.OnOpenStoreBank(sb.owner, npcID)
		case shop != nil && shop.GetCharacterId() == sb.owner.GetID():
			sb.owner.Listener.OnStoreBankLocation(sb.owner, npcID, shop.GetMapId(), uint8(shop.GetChannelId()))
		default:
			sb.owner.Listener.OnStoreBankLocation(sb.owner, npcID, StoreBankNothingMapID, StoreBankNothingChannel)
		}
		return nil
	}).OnError(func(err error) {
		sb.pending = false
		log.Printf("StoreBank.Open character=%d: %v", sb.owner.GetID(), err)
	})
}

func (sb *StoreBank) load(shops []*internal.Shop) {
	sb.shops = shops
	sb.closedAt = time.UnixMilli(shops[0].GetClosedAtUnixMs())
	sb.Meso = 0
	sb.Items = nil
	for _, shop := range shops {
		if closedAt := time.UnixMilli(shop.GetClosedAtUnixMs()); closedAt.Before(sb.closedAt) {
			sb.closedAt = closedAt
		}
		sb.Meso = int32(min(int64(sb.Meso)+int64(shop.GetMeso()), math.MaxInt32))
		for _, pb := range shop.GetItems() {
			item, err := NewItemFromInternalProto(pb.GetItem(), sb.owner.GameWorld)
			if err != nil || pb.GetBundles() == 0 {
				continue
			}
			sb.Items = append(sb.Items, &ShopItem{
				Item:      item,
				Bundles:   uint16(pb.GetBundles()),
				PerBundle: uint16(pb.GetPerBundle()),
				Price:     pb.GetPrice(),
			})
		}
	}
}

func (sb *StoreBank) Close() {
	sb.shops = nil
	sb.Items = nil
	sb.Meso = 0
}

func (sb *StoreBank) fee() (uint32, int32) {
	days := int64(clock.Now().Sub(sb.closedAt) / StoreBankFeeGrace)
	if days <= 0 {
		return 0, 0
	}
	base := int64(sb.Meso)
	for _, listed := range sb.Items {
		base += int64(listed.Price) * int64(listed.Bundles)
	}
	fee := base * min(days, StoreBankFeePercentMax) / 100
	return uint32(days), int32(min(fee, math.MaxInt32))
}

func (sb *StoreBank) Withdraw() error {
	if sb.shops == nil {
		return ErrStoreBankClosed
	}

	days, fee := sb.fee()
	sb.owner.Listener.OnStoreBankFee(sb.owner, days, fee)
	return nil
}

func (sb *StoreBank) check(fee int32) pconst.StoreBankResult {
	if sb.owner.Inventory.Meso < fee {
		return pconst.StoreBankResultNotEnoughMeso
	}
	if int64(sb.owner.Inventory.Meso)-int64(fee)+int64(sb.Meso) > math.MaxInt32 {
		return pconst.StoreBankResultMesoOver
	}
	rewards := map[uint32]uint16{}
	for _, listed := range sb.Items {
		model := listed.Item.GetModel()
		if model.IsOnly() && sb.owner.Inventory.HasItem(model.GetID()) {
			return pconst.StoreBankResultOnlyOne
		}
		rewards[model.GetID()] += listed.count()
	}
	if (ExchangeSpec{Reward: ExchangeSide{Items: rewards}}).Valid(sb.owner) != ExchangeOK {
		return pconst.StoreBankResultInventoryFull
	}
	return pconst.StoreBankResultClaimed
}

func (sb *StoreBank) Confirm(actx actor.Context) error {
	if sb.shops == nil {
		return ErrStoreBankClosed
	}
	if sb.pending {
		return ErrStoreBankBusy
	}
	_, fee := sb.fee()
	if result := sb.check(fee); result != pconst.StoreBankResultClaimed {
		sb.owner.Listener.OnStoreBankResult(sb.owner, result)
		return nil
	}

	sb.pending = true
	sb.owner.Listener.ClaimStoreBankAsync(actx, sb.owner).Do(func(v *internal.ClaimStoreBankReply) error {
		sb.pending = false
		claimed := v.GetShops()
		if len(claimed) == 0 {
			sb.Close()
			sb.owner.Listener.OnStoreBankResult(sb.owner, pconst.StoreBankResultInventoryFull)
			return nil
		}

		sb.load(claimed)
		_, fee := sb.fee()
		if result := sb.check(fee); result != pconst.StoreBankResultClaimed {
			for _, shop := range claimed {
				sb.owner.Listener.OpenShopAsync(actx, sb.owner, shop).OnError(func(err error) {
					log.Printf("StoreBank.Confirm restore character=%d shop=%d: %v", sb.owner.GetID(), shop.GetShopId(), err)
				})
			}
			sb.Close()
			sb.owner.Listener.OnStoreBankResult(sb.owner, result)
			return nil
		}

		sb.owner.Inventory.removeMesoUnchecked(fee)
		sb.owner.Inventory.addMesoUnchecked(sb.Meso)
		for _, listed := range sb.Items {
			sb.owner.Inventory.addItemUnchecked(listed.Item.Clone(listed.count()), true)
		}
		sb.Close()
		sb.owner.Listener.OnStoreBankResult(sb.owner, pconst.StoreBankResultClaimed)
		sb.owner.GameWorld.SaveAsync(actx, []*internal.CharacterSaveEntry{sb.owner.ToProto(sb.owner.GameWorld.GetWorldID())}).OnError(func(err error) {
			log.Printf("StoreBank.Confirm save character=%d: %v", sb.owner.GetID(), err)
		})
		return nil
	}).OnError(func(err error) {
		sb.pending = false
		sb.owner.Listener.OnStoreBankResult(sb.owner, pconst.StoreBankResultInventoryFull)
		log.Printf("StoreBank.Confirm character=%d: %v", sb.owner.GetID(), err)
	})
	return nil
}
