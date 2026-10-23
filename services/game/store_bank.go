package game

import (
	"errors"
	"log"
	"math"
	"slices"
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

var ErrStoreBankClosed = errors.New("store bank is not open")

type StoreBank struct {
	owner    *Character
	shopIDs  []uint32
	claimed  []uint32
	closedAt time.Time
	meso     int64
	items    []*ShopItem
	pending  bool
}

func (sb *StoreBank) Open(actx actor.Context, npcID uint32) {
	if sb.pending {
		return
	}

	sb.pending = true
	sb.owner.gameWorld.Systems().Get[*InternalSystem]().FindEntrustedShopAsync(actx, sb.owner).Bind(sb.owner).Do(func(v *internal.FindEntrustedShopReply) error {
		sb.pending = false
		var shops []*internal.Shop
		for _, closed := range sb.unclaimed(v.GetStoreBank()) {
			if closed.GetCharacterId() == sb.owner.GetID() {
				shops = append(shops, closed)
			}
		}
		shop := v.GetShop()
		switch {
		case len(shops) > 0:
			sb.load(shops)
			sb.owner.listener.OnOpenStoreBank(sb.owner, npcID)
		case shop != nil && shop.GetCharacterId() == sb.owner.GetID():
			sb.owner.listener.OnStoreBankLocation(sb.owner, npcID, shop.GetMapId(), uint8(shop.GetChannelId()))
		default:
			sb.owner.listener.OnStoreBankLocation(sb.owner, npcID, StoreBankNothingMapID, StoreBankNothingChannel)
		}
		return nil
	}).OnError(func(err error) {
		sb.pending = false
		log.Printf("StoreBank.Open character=%d: %v", sb.owner.GetID(), err)
	})
}

func (sb *StoreBank) unclaimed(shops []*internal.Shop) []*internal.Shop {
	return slices.DeleteFunc(slices.Clone(shops), func(shop *internal.Shop) bool {
		return slices.Contains(sb.claimed, shop.GetShopId())
	})
}

func (sb *StoreBank) load(shops []*internal.Shop) {
	sb.shopIDs = make([]uint32, 0, len(shops))
	sb.closedAt = time.UnixMilli(shops[0].GetClosedAtUnixMs())
	sb.meso = 0
	sb.items = nil
	for _, shop := range shops {
		sb.shopIDs = append(sb.shopIDs, shop.GetShopId())
		if closedAt := time.UnixMilli(shop.GetClosedAtUnixMs()); closedAt.Before(sb.closedAt) {
			sb.closedAt = closedAt
		}
		sb.meso += int64(shop.GetMeso())
		for _, pb := range shop.GetItems() {
			item, err := NewItemFromInternalProto(pb.GetItem(), sb.owner.gameWorld)
			if err != nil || pb.GetBundles() == 0 {
				continue
			}
			sb.items = append(sb.items, &ShopItem{
				item:      item,
				bundles:   uint16(pb.GetBundles()),
				perBundle: uint16(pb.GetPerBundle()),
				price:     pb.GetPrice(),
			})
		}
	}
}

func (sb *StoreBank) Close() {
	sb.shopIDs = nil
	sb.items = nil
	sb.meso = 0
}

func (sb *StoreBank) fee() (uint32, int32) {
	days := int64(clock.Now().Sub(sb.closedAt) / StoreBankFeeGrace)
	if days <= 0 {
		return 0, 0
	}
	base := sb.meso
	for _, listed := range sb.items {
		base += int64(listed.item.GetModel().GetPrice()) * int64(listed.count())
	}
	fee := base * min(days, StoreBankFeePercentMax) / 100
	return uint32(days), int32(min(fee, math.MaxInt32))
}

func (sb *StoreBank) Withdraw() error {
	if sb.shopIDs == nil {
		return ErrStoreBankClosed
	}

	days, fee := sb.fee()
	sb.owner.listener.OnStoreBankFee(sb.owner, days, fee)
	return nil
}

func (sb *StoreBank) check(fee int32) pconst.StoreBankResult {
	if sb.owner.inventory.meso < fee {
		return pconst.StoreBankResultNotEnoughMeso
	}
	if int64(sb.owner.inventory.meso)-int64(fee)+sb.meso > math.MaxInt32 {
		return pconst.StoreBankResultMesoOver
	}
	var reward ExchangeSide
	for _, listed := range sb.items {
		if sb.owner.inventory.hasUniqueItem(listed.item.GetModel()) {
			return pconst.StoreBankResultOnlyOne
		}
		reward.include(listed.item, listed.count())
	}
	if (ExchangeSpec{Reward: reward}).Valid(sb.owner, sb.owner.escrow) != ExchangeOK {
		return pconst.StoreBankResultInventoryFull
	}
	return pconst.StoreBankResultClaimed
}

func (sb *StoreBank) Confirm(actx actor.Context) error {
	if sb.shopIDs == nil {
		return ErrStoreBankClosed
	}
	_, fee := sb.fee()
	if result := sb.check(fee); result != pconst.StoreBankResultClaimed {
		sb.owner.listener.OnStoreBankResult(sb.owner, result)
		return nil
	}

	sb.owner.inventory.removeMesoUnchecked(fee)
	sb.owner.inventory.addMesoUnchecked(int32(sb.meso))
	for _, listed := range sb.items {
		sb.owner.inventory.addItemUnchecked(listed.item.Clone(listed.count()), true)
	}
	sb.claimed = append(sb.claimed, sb.shopIDs...)
	sb.Close()
	sb.owner.listener.OnStoreBankResult(sb.owner, pconst.StoreBankResultClaimed)

	entry := sb.owner.BuildRPC(sb.owner.gameWorld.GetWorldID())
	sb.owner.gameWorld.Systems().Get[*InternalSystem]().SaveAsync(actx, []*internal.CharacterSaveEntry{entry}).Bind(sb.owner).Do(func() error {
		sb.claimed = slices.DeleteFunc(sb.claimed, func(id uint32) bool {
			return slices.Contains(entry.GetClaimedShops(), id)
		})
		return nil
	}).OnError(func(err error) {
		log.Printf("StoreBank.Confirm save character=%d: %v", sb.owner.GetID(), err)
	})
	return nil
}

func (sb *StoreBank) CheckShopOpen(actx actor.Context) {
	m := sb.owner.GetMap()
	if m == nil || !m.wz.EntrustedShop || sb.owner.miniRoom != nil {
		sb.owner.listener.OnEntrustedShopCheck(sb.owner, pconst.EntrustedShopCannotOpen, 0, 0)
		return
	}

	sb.owner.gameWorld.Systems().Get[*InternalSystem]().FindEntrustedShopAsync(actx, sb.owner).Bind(sb.owner).Do(func(v *internal.FindEntrustedShopReply) error {
		shop := v.GetShop()
		switch {
		case len(sb.unclaimed(v.GetStoreBank())) > 0:
			sb.owner.listener.OnEntrustedShopCheck(sb.owner, pconst.EntrustedShopStoreBankFull, 0, 0)
		case shop == nil:
			sb.owner.listener.OnEntrustedShopCheck(sb.owner, pconst.EntrustedShopTitle, 0, 0)
		case shop.GetCharacterId() == sb.owner.GetID():
			sb.owner.listener.OnEntrustedShopCheck(sb.owner, pconst.EntrustedShopAlreadyOpen, shop.GetMapId(), uint8(shop.GetChannelId()))
		default:
			sb.owner.listener.OnEntrustedShopCheck(sb.owner, pconst.EntrustedShopAccountBusy, 0, 0)
		}
		return nil
	}).OnError(func(err error) {
		log.Printf("CheckShopOpen character=%d: %v", sb.owner.GetID(), err)
		sb.owner.listener.OnEntrustedShopCheck(sb.owner, pconst.EntrustedShopCannotOpen, 0, 0)

	})
}

func (sb *StoreBank) settleShopOpen(actx actor.Context, request *UnsettledRequest) {
	sb.owner.gameWorld.Systems().Get[*InternalSystem]().AbortOpenShopAsync(actx, sb.owner, request.id).Bind(sb.owner).Do(func(reply *internal.AbortOpenShopReply) error {
		if !reply.GetAborted() {
			request.inFlight = false
			log.Printf("AbortOpenShop character=%d request=%d: shop is no longer empty", sb.owner.GetID(), request.id)
			return nil
		}
		sb.owner.unsettled.Erase(request.id)
		return nil
	}).OnError(func(err error) {
		request.inFlight = false

		log.Printf("AbortOpenShop character=%d request=%d: %v", sb.owner.GetID(), request.id, err)
	})
}
