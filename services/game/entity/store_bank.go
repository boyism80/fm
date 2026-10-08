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
	merchant *internal.HiredMerchant
	Meso     int32
	Items    []*HiredMerchantItem
	pending  bool
}

func (sb *StoreBank) Open(actx actor.Context, npcID uint32) {
	if sb.pending {
		return
	}

	sb.pending = true
	sb.owner.Listener.FindHiredMerchantAsync(actx, sb.owner).Do(func(v *internal.FindHiredMerchantReply) error {
		sb.pending = false
		merchant := v.GetMerchant()
		switch {
		case merchant == nil || merchant.GetCharacterId() != sb.owner.GetID():
			sb.owner.Listener.OnStoreBankLocation(sb.owner, npcID, StoreBankNothingMapID, StoreBankNothingChannel)
		case merchant.GetClosedAtUnixMs() == 0:
			sb.owner.Listener.OnStoreBankLocation(sb.owner, npcID, merchant.GetMapId(), uint8(merchant.GetChannelId()))
		default:
			sb.load(merchant)
			sb.owner.Listener.OnOpenStoreBank(sb.owner, npcID)
		}
		return nil
	}).OnError(func(err error) {
		sb.pending = false
		log.Printf("StoreBank.Open character=%d: %v", sb.owner.GetID(), err)
	})
}

func (sb *StoreBank) load(merchant *internal.HiredMerchant) {
	sb.merchant = merchant
	sb.Meso = merchant.GetMeso()
	sb.Items = nil
	for _, pb := range merchant.GetItems() {
		item, err := NewItemFromInternalProto(pb.GetItem(), sb.owner.GameWorld)
		if err != nil || pb.GetBundles() == 0 {
			continue
		}
		sb.Items = append(sb.Items, &HiredMerchantItem{
			Item:      item,
			Bundles:   uint16(pb.GetBundles()),
			PerBundle: uint16(pb.GetPerBundle()),
			Price:     pb.GetPrice(),
		})
	}
}

func (sb *StoreBank) Close() {
	sb.merchant = nil
	sb.Items = nil
	sb.Meso = 0
}

func (sb *StoreBank) fee() (uint32, int32) {
	days := int64(clock.Now().Sub(time.UnixMilli(sb.merchant.GetClosedAtUnixMs())) / StoreBankFeeGrace)
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
	if sb.merchant == nil {
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
	if sb.merchant == nil {
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
		claimed := v.GetMerchant()
		if claimed == nil {
			sb.Close()
			sb.owner.Listener.OnStoreBankResult(sb.owner, pconst.StoreBankResultInventoryFull)
			return nil
		}

		sb.load(claimed)
		_, fee := sb.fee()
		if result := sb.check(fee); result != pconst.StoreBankResultClaimed {
			sb.owner.Listener.OpenHiredMerchantAsync(actx, sb.owner, claimed).OnError(func(err error) {
				log.Printf("StoreBank.Confirm restore character=%d merchant=%d: %v", sb.owner.GetID(), claimed.GetMerchantId(), err)
			})
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
