package entity

import (
	"errors"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
)

const (
	FreeMarketEntranceMapID = 910000000
	FreeMarketLastMapID     = 910000022
	shopScanRetryDelay      = 5 * time.Second
)

type shopScan struct {
	pending bool
	retryAt time.Time
	found   map[remoteShop]bool
}

func (ch *Character) FindPopularShopSearches(actx actor.Context) {
	ch.Listener.FindPopularShopSearchesAsync(actx, ch).Do(func(v *internal.FindPopularShopSearchesReply) error {
		ch.Listener.OnShopScannerPopular(ch, v.GetItemIds())
		return nil
	}).OnError(func(err error) {
		log.Printf("FindPopularShopSearches character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) UseShopScanner(actx actor.Context, invType constant.InventoryType, slot int16, itemID uint32, searchID uint32, highFirst bool) {
	m := ch.GetMap()
	item := ch.Inventory.GetItem(invType, slot)
	if m == nil || m.TemplateID() < FreeMarketEntranceMapID || m.TemplateID() > FreeMarketLastMapID {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}
	if item == nil || item.GetModel().GetID() != itemID || !constant.IsShopScanner(itemID) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}
	if expiration := item.GetExpiration(); !expiration.IsZero() && clock.Now().After(expiration) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}
	if ch.shopScan.pending || clock.Now().Before(ch.shopScan.retryAt) {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}
	if _, ok := ch.GameWorld.GetResources().Items[searchID]; !ok {
		ch.shopScan.retryAt = clock.Now().Add(shopScanRetryDelay)
		ch.Listener.OnShopScannerResult(ch, searchID, nil)
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}

	ch.shopScan.pending = true
	ch.Listener.SearchShopsAsync(actx, ch, searchID, highFirst).Do(func(v *internal.SearchShopsReply) error {
		ch.shopScan.pending = false
		entries := v.GetEntries()
		ch.shopScan.found = make(map[remoteShop]bool, len(entries))
		for _, entry := range entries {
			ch.shopScan.found[remoteShop{sn: entry.GetSn(), mapID: entry.GetMapId()}] = true
		}
		ch.Listener.OnShopScannerResult(ch, searchID, entries)
		if len(entries) == 0 {
			ch.shopScan.retryAt = clock.Now().Add(shopScanRetryDelay)
		} else if item := ch.Inventory.GetItem(invType, slot); item != nil && item.GetModel().GetID() == itemID {
			ch.Inventory.RemoveItem(invType, slot, 1)
		}
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}).OnError(func(err error) {
		ch.shopScan.pending = false
		ch.Listener.OnUpdateStats(ch, nil, true)
		log.Printf("UseShopScanner character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) VisitShopBySearch(actx actor.Context, sn uint32, mapID uint32) error {
	m := ch.GetMap()
	if m == nil || m.TemplateID() < FreeMarketEntranceMapID || m.TemplateID() > FreeMarketLastMapID {
		return ErrMiniRoomInvalid
	}
	if !ch.shopScan.found[remoteShop{sn: sn, mapID: mapID}] {
		return ErrMiniRoomInvalid
	}
	if m.TemplateID() == mapID {
		return ch.VisitMiniRoom(sn, "")
	}

	target := ch.GameWorld.GetMapSystem().Get(mapID)
	if target == nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	return ch.GameWorld.GetMapSystem().Warp(actx, ch, target, 0, func(actor.Context) {
		ch.RejectMiniRoom(ch.VisitMiniRoom(sn, ""))
	})
}

func (ch *Character) RejectMiniRoom(err error) {
	var enterErr *MiniRoomEnterError
	switch {
	case errors.As(err, &enterErr):
		ch.Listener.OnMiniRoomEnterFailed(ch, enterErr.Code)
	case err != nil:
		ch.Listener.OnUnlockAction(ch)
	}
}
