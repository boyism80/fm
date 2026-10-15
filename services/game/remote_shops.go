package game

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

const (
	shopScanRetryDelay = 5 * time.Second
)

type remoteShop struct {
	sn    uint32
	mapID uint32
}

type RemoteShops struct {
	owner    *Character
	managed  remoteShop
	scanning bool
	retryAt  time.Time
	found    map[remoteShop]bool
}

func (r *RemoteShops) FindPopularSearches(actx actor.Context) {
	r.owner.GameWorld.Systems().Get[*InternalSystem]().FindPopularShopSearchesAsync(actx, r.owner).Do(func(v *internal.FindPopularShopSearchesReply) error {
		r.owner.Listener.OnShopScannerPopular(r.owner, v.GetItemIds())
		return nil
	}).OnError(func(err error) {
		log.Printf("FindPopularShopSearches character=%d: %v", r.owner.GetID(), err)
	})
}

func (r *RemoteShops) Scan(actx actor.Context, invType constant.InventoryType, slot int16, itemID uint32, searchID uint32, highFirst bool) {
	m := r.owner.GetMap()
	item := r.owner.Inventory.GetItem(invType, slot)
	if m == nil || !m.Wz.FreeMarket {
		r.owner.Listener.OnUpdateStats(r.owner, nil, true)
		return
	}
	if item == nil || item.GetModel().GetID() != itemID || item.GetModel().GetTraits()&constant.ItemTraitShopScanner == 0 {
		r.owner.Listener.OnUpdateStats(r.owner, nil, true)
		return
	}
	if expiration := item.GetExpiration(); !expiration.IsZero() && clock.Now().After(expiration) {
		r.owner.Listener.OnUpdateStats(r.owner, nil, true)
		return
	}
	if r.scanning || clock.Now().Before(r.retryAt) {
		r.owner.Listener.OnUpdateStats(r.owner, nil, true)
		return
	}
	if _, ok := r.owner.GameWorld.GetResources().Items[searchID]; !ok {
		r.retryAt = clock.Now().Add(shopScanRetryDelay)
		r.owner.Listener.OnShopScannerResult(r.owner, searchID, nil)
		r.owner.Listener.OnUpdateStats(r.owner, nil, true)
		return
	}

	r.scanning = true
	r.owner.GameWorld.Systems().Get[*InternalSystem]().SearchShopsAsync(actx, r.owner, searchID, highFirst).Do(func(v *internal.SearchShopsReply) error {
		r.scanning = false
		entries := v.GetEntries()
		r.found = make(map[remoteShop]bool, len(entries))
		for _, entry := range entries {
			r.found[remoteShop{sn: entry.GetSn(), mapID: entry.GetMapId()}] = true
		}
		r.owner.Listener.OnShopScannerResult(r.owner, searchID, entries)
		if len(entries) == 0 {
			r.retryAt = clock.Now().Add(shopScanRetryDelay)
		} else if item := r.owner.Inventory.GetItem(invType, slot); item != nil && item.GetModel().GetID() == itemID {
			r.owner.Inventory.RemoveItem(invType, slot, 1)
		}
		r.owner.Listener.OnUpdateStats(r.owner, nil, true)
		return nil
	}).OnError(func(err error) {
		r.scanning = false
		r.owner.Listener.OnUpdateStats(r.owner, nil, true)
		log.Printf("UseShopScanner character=%d: %v", r.owner.GetID(), err)
	})
}

func (r *RemoteShops) VisitFound(actx actor.Context, sn uint32, mapID uint32) error {
	m := r.owner.GetMap()
	if m == nil || !m.Wz.FreeMarket {
		return ErrMiniRoomInvalid
	}
	if !r.found[remoteShop{sn: sn, mapID: mapID}] {
		return ErrMiniRoomInvalid
	}
	if m.TemplateID() == mapID {
		return m.VisitMiniRoom(r.owner, sn, "")
	}

	target := r.owner.GameWorld.Systems().Get[*MapSystem]().Get(mapID)
	if target == nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	return r.owner.GameWorld.Systems().Get[*MapSystem]().Warp(actx, r.owner, target, 0, func(actor.Context) {
		r.owner.Listener.OnMiniRoomFailed(r.owner, target.VisitMiniRoom(r.owner, sn, ""))
	})
}

func (r *RemoteShops) Manage(actx actor.Context, slot int16) error {
	item := r.owner.Inventory.GetItem(constant.InventoryTypeCash, slot)
	if item == nil || item.GetModel().GetID() != RemoteEntrustedShopItemID {
		return ErrMiniRoomInvalid
	}
	if expiration := item.GetExpiration(); !expiration.IsZero() && clock.Now().After(expiration) {
		return ErrMiniRoomInvalid
	}
	if r.owner.MiniRoom != nil {
		return ErrMiniRoomBusy
	}

	r.owner.GameWorld.Systems().Get[*InternalSystem]().FindEntrustedShopAsync(actx, r.owner).Do(func(v *internal.FindEntrustedShopReply) error {
		r.owner.Call(func(actor.Context) {
			shop := v.GetShop()
			target := r.owner.GameWorld.Systems().Get[*MapSystem]().Get(shop.GetMapId())
			switch {
			case shop == nil || shop.GetCharacterId() != r.owner.GetID():
				r.owner.Listener.OnEntrustedShopCheck(r.owner, pconst.EntrustedShopRemoteLocation, StoreBankNothingMapID, StoreBankNothingChannel)
			case shop.GetChannelId() != r.owner.GameWorld.GetChannelID():
				r.owner.Listener.OnEntrustedShopCheck(r.owner, pconst.EntrustedShopRemoteLocation, shop.GetMapId(), uint8(shop.GetChannelId()))
			case target == nil:
				r.owner.Listener.OnEntrustedShopCheck(r.owner, pconst.EntrustedShopRemoteLocation, StoreBankNothingMapID, StoreBankNothingChannel)
			default:
				r.owner.GameWorld.Systems().Get[*MapSystem]().Call(target, func(actor.Context) {
					sn, ok := target.FindEntrustedShopByOwner(r.owner.GetID())
					r.owner.GameWorld.Systems().Get[*DispatchSystem]().CallCharacter(r.owner.GetID(), func(ctx actor.Context, c *Character) {
						if c == nil {
							return
						}
						if !ok {
							c.Listener.OnEntrustedShopCheck(c, pconst.EntrustedShopRemoteLocation, StoreBankNothingMapID, StoreBankNothingChannel)
							return
						}
						c.RemoteShops.managed = remoteShop{sn: sn, mapID: target.TemplateID()}
						c.Listener.OnEntrustedShopRemoteVisit(c, sn)
					})
				})
			}
		})
		return nil
	}).OnError(func(err error) {
		log.Printf("UseRemoteEntrustedShop character=%d: %v", r.owner.GetID(), err)
		r.owner.Call(func(actor.Context) {
			r.owner.Listener.OnUnlockAction(r.owner)
		})
	})
	return nil
}

func (r *RemoteShops) visitManaged(sn uint32) error {
	if r.managed.sn != sn {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	if r.owner.MiniRoom != nil {
		return ErrMiniRoomInvalid
	}

	managed := r.managed
	r.managed = remoteShop{}
	target := r.owner.GameWorld.Systems().Get[*MapSystem]().Get(managed.mapID)
	if target == nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	ownerID := r.owner.GetID()
	r.owner.GameWorld.Systems().Get[*MapSystem]().Call(target, func(actor.Context) {
		es, ok := target.GetObject(constant.ObjectTypeEntrustedShop, sn).(*EntrustedShop)
		if !ok {
			r.owner.GameWorld.Systems().Get[*DispatchSystem]().CallCharacter(ownerID, func(ctx actor.Context, c *Character) {
				if c == nil {
					return
				}
				c.Listener.OnMiniRoomEnterFailed(c, pconst.MiniRoomEnterClosed)
			})
			return
		}
		es.Manage(ownerID)
	})
	return nil
}
