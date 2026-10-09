package entity

import (
	"fmt"
	"log"
	"math"
	"slices"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

const (
	EntrustedShopMaxItems     = 10
	EntrustedShopDuration     = 24 * time.Hour
	RemoteEntrustedShopItemID = 5470000
	entrustedShopCloseTimer   = "entrusted_shop_close"
	entrustedShopKeptMessage  = "인벤토리에 공간이 부족해 받지 못한 물품은 프레드릭에게 보관되었습니다."
)

type EntrustedShop struct {
	shopRoom
	soldInform bool
	managing   bool
}

type EntrustedShopView struct {
	ItemID    uint32
	OwnerName string
	Title     string
	Elapsed   time.Duration
	Meso      int32
	Items     []*ShopItem
	Sold      []ShopSale
	Visitors  [ShopVisitors]*dto.Character
}

func (es *EntrustedShop) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeEntrustedShop
}

func (es *EntrustedShop) Is(typ constant.ObjectType) bool {
	return es.GetObjectType().Has(typ)
}

func (es *EntrustedShop) view() EntrustedShopView {
	view := EntrustedShopView{
		ItemID:    es.ItemID,
		OwnerName: es.OwnerName,
		Title:     es.Title,
		Elapsed:   es.Elapsed(),
		Meso:      es.Meso,
		Sold:      slices.Clone(es.Sold),
	}
	for _, listed := range es.Items {
		copied := *listed
		view.Items = append(view.Items, &copied)
	}
	for i, visitor := range es.Visitors {
		if visitor != nil {
			view.Visitors[i] = visitor.ToDTO()
		}
	}
	return view
}

func (es *EntrustedShop) balloon() response.EntrustedShopBalloon {
	return response.EntrustedShopBalloon{
		SN:     es.OID,
		Title:  es.Title,
		ItemID: es.ItemID,
		Users:  es.users(),
	}
}

func (es *EntrustedShop) SendSpawnSyncToViewer(viewer *Character) {
	if es.published == false {
		return
	}
	footholdID := uint16(0)
	if foothold, ok := es.Map.Wz.Footholds.Find(types.Point[int16]{X: es.Position.X, Y: es.Position.Y}); ok {
		footholdID = uint16(foothold.ID)
	}
	viewer.Send(&response.SpawnEntrustedShop{
		EmployerID:           es.OwnerID,
		X:                    es.Position.X,
		Y:                    es.Position.Y,
		Foothold:             footholdID,
		OwnerName:            es.OwnerName,
		EntrustedShopBalloon: es.balloon(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (es *EntrustedShop) SendDestroySyncToViewer(viewer *Character) {
	if es.published == false {
		return
	}
	viewer.Send(&response.DestroyEntrustedShop{EmployerID: es.OwnerID}, types.SEND_POLICY_ENCRYPT)
}

func (es *EntrustedShop) updateBalloon() {
	if es.published == false || es.Map == nil {
		return
	}
	es.Broadcast(&response.UpdateEntrustedShop{
		EmployerID:           es.OwnerID,
		EntrustedShopBalloon: es.balloon(),
	}, nil)
}

func (es *EntrustedShop) accepting() bool {
	return es.published && es.managing == false && es.Map != nil
}

func (es *EntrustedShop) release(member *Character, slot uint8, reason pconst.MiniRoomLeaveReason) {
	es.GameWorld.GetDispatchSystem().CallCharacter(member.GetID(), func(ctx actor.Context, c *Character) {
		if c == nil || c.MiniRoom != es {
			return
		}
		c.MiniRoom = nil
		c.Listener.OnMiniRoomLeft(c, slot, reason)
	})
}

func (es *EntrustedShop) notifyOwner(notify func(owner *Character)) {
	es.GameWorld.GetDispatchSystem().CallCharacter(es.OwnerID, func(ctx actor.Context, c *Character) {
		if c == nil || c.MiniRoom != es {
			return
		}
		notify(c)
	})
}

func (es *EntrustedShop) payOwner(meso int32, items []*ShopItem, closing bool, notify func(owner *Character)) {
	worldID := es.GameWorld.GetWorldID()
	es.GameWorld.GetDispatchSystem().CallCharacter(es.OwnerID, func(ctx actor.Context, c *Character) {
		var entry *internal.CharacterSaveEntry
		if c != nil {
			meso, items = c.collectShopGoods(meso, items)
			entry = c.ToProto(worldID)
			if c.MiniRoom == es {
				notify(c)
			}
			if meso > 0 || len(items) > 0 {
				c.Listener.OnMessage(c, constant.MsgPopup, entrustedShopKeptMessage)
			}
		}

		es.GameWorld.GetMapSystem().Call(es.home, func(ctx actor.Context) {
			if entry != nil {
				es.entries[es.OwnerID] = entry
			}
			if meso > 0 || len(items) > 0 {
				es.storeBank = append(es.storeBank, es.storeBankToProto(meso, items))
			}
			es.save(ctx, closing)
		})
	})
}

func (es *EntrustedShop) Elapsed() time.Duration {
	return clock.Now().Sub(es.OpenedAt)
}

func (es *EntrustedShop) Visit(ch *Character) error {
	if ch.MiniRoom != nil {
		return ErrMiniRoomInvalid
	}
	if es.published == false || es.Map == nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	if es.managing {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterOrganizing}
	}
	index, ok := es.freeSlot()
	if ok == false {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterFull}
	}

	slot := uint8(index + 1)
	for _, member := range es.Members() {
		member.Listener.OnMiniRoomVisited(member, slot, ch)
	}
	es.Visitors[index] = ch
	ch.MiniRoom = es
	es.updateBalloon()
	ch.Listener.OnMiniRoomEntered(ch, slot, es.view(), false)
	return nil
}

func (es *EntrustedShop) Manage(ownerID uint32) {
	if es.accepting() == false || es.OwnerID != ownerID {
		es.GameWorld.GetDispatchSystem().CallCharacter(ownerID, func(ctx actor.Context, c *Character) {
			if c == nil {
				return
			}
			c.Listener.OnMiniRoomEnterFailed(c, pconst.MiniRoomEnterClosed)
		})
		return
	}

	for i, visitor := range es.Visitors {
		if visitor == nil {
			continue
		}
		es.Visitors[i] = nil
		es.release(visitor, uint8(i+1), pconst.MiniRoomLeaveOrganizing)
	}
	es.managing = true
	es.updateBalloon()

	view := es.view()
	es.GameWorld.GetDispatchSystem().CallCharacter(ownerID, func(ctx actor.Context, c *Character) {
		if c != nil && c.MiniRoom == nil {
			c.MiniRoom = es
			c.Listener.OnMiniRoomEntered(c, 0, view, false)
			return
		}
		es.GameWorld.GetMapSystem().Call(es.home, es.stopManaging)
	})
}

func (es *EntrustedShop) stopManaging(ctx actor.Context) {
	if es.managing == false {
		return
	}

	es.managing = false
	if es.published {
		es.updateBalloon()
		return
	}
	es.close(ctx)
}

func (es *EntrustedShop) Leave(ch *Character) {
	if ch.GetID() == es.OwnerID {
		ch.MiniRoom = nil
		es.GameWorld.GetMapSystem().Call(es.home, es.stopManaging)
		return
	}

	slot, ok := es.SlotOf(ch)
	if ok == false {
		return
	}
	ch.MiniRoom = nil
	es.Visitors[slot-1] = nil
	for _, member := range es.Members() {
		member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
	}
	es.updateBalloon()
}

func (es *EntrustedShop) Chat(ch *Character, message string) error {
	if ch.GetID() == es.OwnerID {
		ch.Listener.OnMiniRoomChat(ch, 0, fmt.Sprintf("%s : %s", ch.GetName(), message))
		return nil
	}
	return es.shopRoom.Chat(ch, message)
}

func (es *EntrustedShop) AddItem(actx actor.Context, ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}
	listed, err := ch.takeShopItem(invType, slot, bundles, perBundle, price)
	if err != nil {
		return err
	}

	entry := ch.ToProto(es.GameWorld.GetWorldID())
	es.GameWorld.GetMapSystem().Call(es.home, func(ctx actor.Context) {
		if es.managing == false || len(es.Items) >= EntrustedShopMaxItems || es.holds(listed.Item.GetModel().GetID()) {
			es.payOwner(0, []*ShopItem{listed}, false, func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		es.Items = append(es.Items, listed)
		es.entries[es.OwnerID] = entry
		es.save(ctx, false)
		view := es.view()
		es.notifyOwner(func(owner *Character) {
			owner.Listener.OnMiniRoomItems(owner, view)
		})
	})
	return nil
}

func (es *EntrustedShop) RemoveItem(actx actor.Context, ch *Character, index uint16) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.GetMapSystem().Call(es.home, func(ctx actor.Context) {
		if es.managing == false || int(index) >= len(es.Items) {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		listed := es.Items[index]
		es.Items = append(es.Items[:index:index], es.Items[index+1:]...)
		var items []*ShopItem
		if listed.Bundles > 0 {
			items = append(items, listed)
		}
		count := uint8(len(es.Items))
		es.payOwner(0, items, false, func(owner *Character) {
			owner.Listener.OnMiniRoomItemRemoved(owner, count, index)
		})
	})
	return nil
}

func (es *EntrustedShop) Open(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.GetMapSystem().Call(es.home, func(ctx actor.Context) {
		if es.managing == false || es.published || len(es.Items) == 0 || es.Map == nil {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		es.managing = false
		es.published = true
		es.BroadcastCall(func(obj Object) {
			viewer, ok := obj.(*Character)
			if ok == false {
				return
			}
			es.SendSpawnSyncToViewer(viewer)
		}, nil)
		es.scheduleClose()
		es.save(ctx, false)
		es.notifyOwner(func(owner *Character) {
			owner.MiniRoom = nil
		})
	})
	return nil
}

func (es *EntrustedShop) scheduleClose() {
	remaining := max(EntrustedShopDuration-es.Elapsed(), 0)
	es.AddTimer(entrustedShopCloseTimer, remaining, false, func() {
		es.GameWorld.GetMapSystem().Call(es.home, func(ctx actor.Context) {
			if es.Map == nil {
				return
			}
			for i, visitor := range es.Visitors {
				if visitor == nil {
					continue
				}
				es.release(visitor, uint8(i+1), pconst.MiniRoomLeaveTimeUp)
			}
			es.Visitors = [ShopVisitors]*Character{}
			if es.managing {
				es.managing = false
				es.notifyOwner(func(owner *Character) {
					owner.MiniRoom = nil
					owner.Listener.OnMiniRoomLeft(owner, 0, pconst.MiniRoomLeaveTimeUp)
				})
			}
			es.close(ctx)
		})
	})
}

func (es *EntrustedShop) EndMaintenance(ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.Leave(ch)
	return nil
}

func (es *EntrustedShop) Buy(actx actor.Context, ch *Character, index uint16, bundles uint16) error {
	if _, ok := es.SlotOf(ch); ok == false || es.accepting() == false {
		return ErrMiniRoomInvalid
	}

	listed, total, income, err := es.sell(ch, index, bundles, es.Meso)
	if err != nil {
		return err
	}
	es.Meso += income
	model := listed.Item.GetModel()
	es.Sold = append(es.Sold, ShopSale{
		ItemID:  model.GetID(),
		Bundles: bundles,
		Total:   total,
		Buyer:   ch.GetName(),
	})
	if len(es.Sold) > math.MaxUint8 {
		es.Sold = es.Sold[len(es.Sold)-math.MaxUint8:]
	}
	view := es.view()
	for _, member := range es.Members() {
		member.Listener.OnMiniRoomItems(member, view)
	}
	es.save(actx, false, ch)

	if es.soldInform == false {
		return nil
	}
	message := fmt.Sprintf("고용상점에서 %s %d개가 판매되었습니다.", es.GameWorld.GetResources().GetItemName(model.GetID()), bundles*listed.PerBundle)
	es.GameWorld.GetDispatchSystem().CallCharacter(es.OwnerID, func(ctx actor.Context, owner *Character) {
		if owner == nil {
			return
		}
		owner.Listener.OnMessage(owner, constant.MsgLightBlueText, message)
	})
	return nil
}

func (es *EntrustedShop) Arrange(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.GetMapSystem().Call(es.home, func(ctx actor.Context) {
		if es.managing == false {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.Meso
		es.Meso = 0
		items := es.Items[:0]
		for _, listed := range es.Items {
			if listed.Bundles > 0 {
				items = append(items, listed)
			}
		}
		es.Items = items
		view := es.view()
		es.payOwner(meso, nil, false, func(owner *Character) {
			owner.Listener.OnMiniRoomArranged(owner, view)
		})
	})
	return nil
}

func (es *EntrustedShop) WithdrawMeso(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.GetMapSystem().Call(es.home, func(ctx actor.Context) {
		if es.managing == false || es.Meso <= 0 {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.Meso
		es.Meso = 0
		es.payOwner(meso, nil, false, func(owner *Character) {
			owner.Listener.OnMiniRoomMesoWithdrawn(owner)
		})
	})
	return nil
}

func (es *EntrustedShop) Close(actx actor.Context, ch *Character) error {
	if ch.GetID() != es.OwnerID {
		return ErrMiniRoomNotOwner
	}

	es.GameWorld.GetMapSystem().Call(es.home, func(ctx actor.Context) {
		if es.managing == false || es.Map == nil {
			es.notifyOwner(func(owner *Character) {
				owner.Listener.OnUnlockAction(owner)
			})
			return
		}

		meso := es.Meso
		var items []*ShopItem
		for _, listed := range es.Items {
			if listed.Bundles > 0 {
				items = append(items, listed)
			}
		}
		es.Meso = 0
		es.Items = nil
		es.managing = false
		es.Map.RemoveEntrustedShop(es)
		es.payOwner(meso, items, true, func(owner *Character) {
			owner.MiniRoom = nil
			owner.Listener.OnMiniRoomClosed(owner, pconst.MiniRoomCloseAll)
		})
	})
	return nil
}

func (es *EntrustedShop) close(actx actor.Context) {
	m := es.Map
	if m == nil {
		return
	}
	m.RemoveEntrustedShop(es)
	es.save(actx, true)
}

func (ch *Character) UseEntrustedShop(actx actor.Context) {
	m := ch.GetMap()
	if m == nil || m.Wz.EntrustedShop == false || ch.MiniRoom != nil || ch.miniRoomPending {
		ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopCannotOpen, 0, 0)
		return
	}

	ch.miniRoomPending = true
	ch.Listener.FindEntrustedShopAsync(actx, ch).Do(func(v *internal.FindEntrustedShopReply) error {
		ch.miniRoomPending = false
		shop := v.GetShop()
		switch {
		case len(v.GetStoreBank()) > 0:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopStoreBankFull, 0, 0)
		case shop == nil:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopTitle, 0, 0)
		case shop.GetCharacterId() == ch.GetID():
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopAlreadyOpen, shop.GetMapId(), uint8(shop.GetChannelId()))
		default:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopAccountBusy, 0, 0)
		}
		return nil
	}).OnError(func(err error) {
		ch.miniRoomPending = false
		ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopCannotOpen, 0, 0)
		log.Printf("UseEntrustedShop character=%d: %v", ch.GetID(), err)
	})
}

func (ch *Character) CreateEntrustedShop(actx actor.Context, title string, slot int16, itemID uint32) error {
	m := ch.GetMap()
	if m == nil || m.Wz.EntrustedShop == false || ch.MiniRoom != nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterCannotOpen}
	}
	if ch.miniRoomPending {
		return ErrMiniRoomBusy
	}
	permit := ch.Inventory.GetItem(constant.InventoryTypeCash, slot)
	if permit == nil || permit.GetModel().GetID() != itemID || itemID/10000 != 503 {
		return ErrMiniRoomInvalid
	}
	if expiration := permit.GetExpiration(); expiration.IsZero() == false && clock.Now().After(expiration) {
		return ErrMiniRoomInvalid
	}
	if title == "" {
		return ErrMiniRoomInvalid
	}
	if err := m.checkShopSpot(ch.Position); err != nil {
		return err
	}

	soldInform := false
	if model, ok := permit.GetModel().(*wz.CashItem); ok {
		soldInform = model.SoldInform
	}
	es := &EntrustedShop{
		shopRoom: shopRoom{
			kind:      pconst.MiniRoomTypeEntrustedShop,
			AccountID: ch.AccountID,
			OwnerID:   ch.GetID(),
			OwnerName: ch.GetName(),
			MapID:     m.TemplateID(),
			ItemID:    itemID,
			Title:     title,
			OpenedAt:  clock.Now(),
			entries:   make(map[uint32]*internal.CharacterSaveEntry),
		},
		soldInform: soldInform,
	}
	es.ObjectCore.self = es
	es.Position = ch.Position

	ch.miniRoomPending = true
	ch.Listener.OpenShopAsync(actx, ch, es.ToProto()).Do(func(v *internal.OpenShopReply) error {
		ch.miniRoomPending = false
		if v.GetExisting() != nil {
			ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
			return nil
		}

		es.ID = v.GetShopId()
		current := ch.GetMap()
		if current != m || ch.MiniRoom != nil || m.checkShopSpot(es.Position) != nil {
			es.GameWorld = ch.GameWorld
			es.save(actx, true)
			ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
			return nil
		}
		m.AddEntrustedShop(es)
		es.managing = true
		ch.MiniRoom = es
		ch.Listener.OnMiniRoomEntered(ch, 0, es.view(), true)
		return nil
	}).OnError(func(err error) {
		ch.miniRoomPending = false
		ch.Listener.OnMiniRoomEnterFailed(ch, pconst.MiniRoomEnterCannotOpen)
		log.Printf("CreateEntrustedShop character=%d: %v", ch.GetID(), err)
	})
	return nil
}

func (ch *Character) UseRemoteEntrustedShop(actx actor.Context, slot int16) error {
	item := ch.Inventory.GetItem(constant.InventoryTypeCash, slot)
	if item == nil || item.GetModel().GetID() != RemoteEntrustedShopItemID {
		return ErrMiniRoomInvalid
	}
	if expiration := item.GetExpiration(); expiration.IsZero() == false && clock.Now().After(expiration) {
		return ErrMiniRoomInvalid
	}
	if ch.MiniRoom != nil || ch.miniRoomPending {
		return ErrMiniRoomBusy
	}

	ch.miniRoomPending = true
	ch.Listener.FindEntrustedShopAsync(actx, ch).Do(func(v *internal.FindEntrustedShopReply) error {
		ch.miniRoomPending = false
		shop := v.GetShop()
		target := ch.GameWorld.GetMapSystem().Get(shop.GetMapId())
		switch {
		case shop == nil || shop.GetCharacterId() != ch.GetID():
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopRemoteLocation, StoreBankNothingMapID, StoreBankNothingChannel)
		case shop.GetChannelId() != ch.GameWorld.GetChannelID():
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopRemoteLocation, shop.GetMapId(), uint8(shop.GetChannelId()))
		case target == nil:
			ch.Listener.OnEntrustedShopCheck(ch, pconst.EntrustedShopRemoteLocation, StoreBankNothingMapID, StoreBankNothingChannel)
		default:
			ch.GameWorld.GetMapSystem().Call(target, func(actor.Context) {
				sn, ok := target.FindEntrustedShopByOwner(ch.GetID())
				ch.GameWorld.GetDispatchSystem().CallCharacter(ch.GetID(), func(ctx actor.Context, c *Character) {
					if c == nil {
						return
					}
					if ok == false {
						c.Listener.OnEntrustedShopCheck(c, pconst.EntrustedShopRemoteLocation, StoreBankNothingMapID, StoreBankNothingChannel)
						return
					}
					c.remoteShop = remoteShop{sn: sn, mapID: target.TemplateID()}
					c.Listener.OnEntrustedShopRemoteVisit(c, sn)
				})
			})
		}
		return nil
	}).OnError(func(err error) {
		ch.miniRoomPending = false
		ch.Listener.OnUnlockAction(ch)
		log.Printf("UseRemoteEntrustedShop character=%d: %v", ch.GetID(), err)
	})
	return nil
}

func (ch *Character) VisitMiniRoom(sn uint32) error {
	m := ch.GetMap()
	if m == nil {
		return ErrMiniRoomInvalid
	}
	if es, ok := m.GetObject(constant.ObjectTypeEntrustedShop, sn).(*EntrustedShop); ok {
		if ch.GetID() != es.OwnerID {
			return es.Visit(ch)
		}
		if ch.MiniRoom != nil {
			return ErrMiniRoomInvalid
		}
		es.Manage(ch.GetID())
		return nil
	}
	if ps, ok := m.GetObject(constant.ObjectTypePersonalShop, sn).(*PersonalShop); ok {
		return ps.Visit(ch)
	}
	if t, ok := m.GetObject(constant.ObjectTypeTrade, sn).(*Trade); ok {
		return t.Visit(ch)
	}
	if ch.remoteShop.sn != sn {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	if ch.MiniRoom != nil {
		return ErrMiniRoomInvalid
	}

	remote := ch.remoteShop
	ch.remoteShop = remoteShop{}
	target := ch.GameWorld.GetMapSystem().Get(remote.mapID)
	if target == nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	ownerID := ch.GetID()
	ch.GameWorld.GetMapSystem().Call(target, func(actor.Context) {
		es, ok := target.GetObject(constant.ObjectTypeEntrustedShop, sn).(*EntrustedShop)
		if ok == false {
			ch.GameWorld.GetDispatchSystem().CallCharacter(ownerID, func(ctx actor.Context, c *Character) {
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
