package entity

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	pconst "github.com/boyism80/fm/protocol/constant"
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
)

type EntrustedShop struct {
	shopRoom
	soldInform bool
	remote     bool
	ended      bool
}

func (es *EntrustedShop) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeEntrustedShop
}

func (es *EntrustedShop) Is(typ constant.ObjectType) bool {
	return es.GetObjectType().Has(typ)
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
	es.mu.Lock()
	defer es.mu.Unlock()
	es.sendSpawn(viewer)
}

func (es *EntrustedShop) sendSpawn(viewer *Character) {
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
	es.mu.Lock()
	defer es.mu.Unlock()
	es.sendDestroy(viewer)
}

func (es *EntrustedShop) sendDestroy(viewer *Character) {
	if es.published == false {
		return
	}
	viewer.Send(&response.DestroyEntrustedShop{EmployerID: es.OwnerID}, types.SEND_POLICY_ENCRYPT)
}

func (es *EntrustedShop) updateBalloon() {
	if es.published == false || es.Map == nil {
		return
	}
	es.GameWorld.GetMapSystem().Call(es.Map, func(actor.Context) {
		es.mu.Lock()
		defer es.mu.Unlock()
		if es.Map == nil {
			return
		}
		es.Broadcast(&response.UpdateEntrustedShop{
			EmployerID:           es.OwnerID,
			EntrustedShopBalloon: es.balloon(),
		}, nil)
	})
}

func (es *EntrustedShop) accepting() bool {
	return es.published && es.owner == nil && es.ended == false
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

func (es *EntrustedShop) Elapsed() time.Duration {
	return clock.Now().Sub(es.OpenedAt)
}

func (es *EntrustedShop) Visit(ch *Character) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if ch.MiniRoom != nil {
		return ErrMiniRoomInvalid
	}

	if ch.GetID() == es.OwnerID {
		if es.accepting() == false || es.Map == nil {
			return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
		}
		for i, visitor := range es.Visitors {
			if visitor == nil {
				continue
			}
			es.Visitors[i] = nil
			es.release(visitor, uint8(i+1), pconst.MiniRoomLeaveOrganizing)
		}
		es.owner = ch
		es.remote = ch.GetMap() != es.Map
		ch.MiniRoom = es
		es.updateBalloon()
		ch.Listener.OnMiniRoomEntered(ch, es, false)
		return nil
	}

	if es.published == false || es.ended {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	if es.owner != nil {
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
	ch.Listener.OnMiniRoomEntered(ch, es, false)
	return nil
}

func (es *EntrustedShop) Leave(ch *Character) {
	es.mu.Lock()
	defer es.mu.Unlock()
	slot, ok := es.SlotOf(ch)
	if ok == false {
		return
	}

	ch.MiniRoom = nil
	if slot != 0 {
		es.Visitors[slot-1] = nil
		for _, member := range es.Members() {
			member.Listener.OnMiniRoomLeft(member, slot, pconst.MiniRoomLeaveExit)
		}
		es.updateBalloon()
		return
	}

	es.owner = nil
	es.remote = false
	if es.published {
		es.updateBalloon()
		return
	}
	es.queueClose()
}

func (es *EntrustedShop) Chat(ch *Character, message string) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	return es.shopRoom.Chat(ch, message)
}

func (es *EntrustedShop) AddItem(actx actor.Context, ch *Character, invType constant.InventoryType, slot int16, bundles uint16, perBundle uint16, price int32) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if len(es.Items) >= EntrustedShopMaxItems {
		return ErrMiniRoomFull
	}

	if err := es.list(ch, invType, slot, bundles, perBundle, price); err != nil {
		return err
	}
	ch.Listener.OnMiniRoomItems(ch, es)
	es.save(actx, false, ch)
	return nil
}

func (es *EntrustedShop) RemoveItem(actx actor.Context, ch *Character, index uint16) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}

	if err := es.unlist(ch, index); err != nil {
		return err
	}
	ch.Listener.OnMiniRoomItems(ch, es)
	es.save(actx, false, ch)
	return nil
}

func (es *EntrustedShop) Open(actx actor.Context, ch *Character) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if es.owner != ch || es.published {
		return ErrMiniRoomNotOwner
	}
	if len(es.Items) == 0 {
		return ErrMiniRoomInvalid
	}

	es.owner = nil
	ch.MiniRoom = nil
	es.published = true
	es.BroadcastCall(func(obj Object) {
		viewer, ok := obj.(*Character)
		if ok == false {
			return
		}
		es.sendSpawn(viewer)
	}, nil)
	es.scheduleClose()
	es.save(actx, false)
	return nil
}

func (es *EntrustedShop) scheduleClose() {
	remaining := max(EntrustedShopDuration-es.Elapsed(), 0)
	es.AddTimer(entrustedShopCloseTimer, remaining, false, func() {
		es.GameWorld.GetMapSystem().Call(es.Map, func(ctx actor.Context) {
			es.mu.Lock()
			defer es.mu.Unlock()
			if es.Map == nil {
				return
			}
			for _, member := range es.Members() {
				slot, _ := es.SlotOf(member)
				es.release(member, slot, pconst.MiniRoomLeaveTimeUp)
			}
			es.owner = nil
			es.remote = false
			es.Visitors = [ShopVisitors]*Character{}
			es.close(ctx)
		})
	})
}

func (es *EntrustedShop) EndMaintenance(ch *Character) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if es.owner != ch || es.published == false {
		return ErrMiniRoomNotOwner
	}

	es.owner = nil
	es.remote = false
	ch.MiniRoom = nil
	es.updateBalloon()
	return nil
}

func (es *EntrustedShop) Buy(actx actor.Context, ch *Character, index uint16, bundles uint16) error {
	es.mu.Lock()
	defer es.mu.Unlock()
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
	for _, member := range es.Members() {
		member.Listener.OnMiniRoomItems(member, es)
	}
	es.save(actx, false, ch)

	if es.soldInform == false {
		return nil
	}
	message := fmt.Sprintf("고용상점에서 %s %d개가 판매되었습니다.", es.GameWorld.GetResources().GetItemName(model.GetID()), bundles*listed.PerBundle)
	es.GameWorld.GetDispatchSystem().CallCharacter(es.OwnerID, func(ctx actor.Context, owner *Character) {
		owner.Listener.OnMessage(owner, constant.MsgLightBlueText, message)
	})
	return nil
}

func (es *EntrustedShop) Arrange(actx actor.Context, ch *Character) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}

	if (ExchangeSpec{Reward: ExchangeSide{Meso: es.Meso}}).Valid(ch) == ExchangeOK {
		ch.Inventory.addMesoUnchecked(es.Meso)
		es.Meso = 0
	}
	items := es.Items[:0]
	for _, listed := range es.Items {
		if listed.Bundles > 0 {
			items = append(items, listed)
		}
	}
	es.Items = items
	ch.Listener.OnMiniRoomArranged(ch, es)
	es.save(actx, false, ch)
	return nil
}

func (es *EntrustedShop) WithdrawMeso(actx actor.Context, ch *Character) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}
	if es.Meso <= 0 {
		return ErrMiniRoomInvalid
	}
	if (ExchangeSpec{Reward: ExchangeSide{Meso: es.Meso}}).Valid(ch) != ExchangeOK {
		return ErrMiniRoomMesoOver
	}

	ch.Inventory.addMesoUnchecked(es.Meso)
	es.Meso = 0
	ch.Listener.OnMiniRoomMesoWithdrawn(ch)
	es.save(actx, false, ch)
	return nil
}

func (es *EntrustedShop) Close(actx actor.Context, ch *Character) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	if es.owner != ch {
		return ErrMiniRoomNotOwner
	}

	result := pconst.MiniRoomCloseAll
	if (ExchangeSpec{Reward: ExchangeSide{Meso: es.Meso}}).Valid(ch) != ExchangeOK {
		result = pconst.MiniRoomCloseMesoOver
	} else {
		ch.Inventory.addMesoUnchecked(es.Meso)
		es.Meso = 0
		result = es.returnItems(ch)
	}

	es.owner = nil
	es.remote = false
	es.ended = true
	ch.MiniRoom = nil
	ch.Listener.OnMiniRoomClosed(ch, result)
	es.entries[ch.GetID()] = ch.ToProto(es.GameWorld.GetWorldID())
	es.queueClose()
	return nil
}

func (es *EntrustedShop) returnItems(ch *Character) pconst.MiniRoomCloseResult {
	rewards := map[uint32]uint16{}
	for _, listed := range es.Items {
		if listed.Bundles == 0 {
			continue
		}
		model := listed.Item.GetModel()
		if model.IsOnly() && ch.Inventory.HasItem(model.GetID()) {
			return pconst.MiniRoomCloseOnlyOne
		}
		rewards[model.GetID()] += listed.count()
	}
	if (ExchangeSpec{Reward: ExchangeSide{Items: rewards}}).Valid(ch) != ExchangeOK {
		return pconst.MiniRoomCloseInventoryFull
	}

	for _, listed := range es.Items {
		if listed.Bundles == 0 {
			continue
		}
		ch.Inventory.addItemUnchecked(listed.Item.Clone(listed.count()), true)
	}
	es.Items = nil
	return pconst.MiniRoomCloseAll
}

func (es *EntrustedShop) queueClose() {
	if es.Map == nil {
		return
	}
	es.GameWorld.GetMapSystem().Call(es.Map, func(ctx actor.Context) {
		es.mu.Lock()
		defer es.mu.Unlock()
		es.close(ctx)
	})
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
		es.owner = ch
		ch.MiniRoom = es
		ch.Listener.OnMiniRoomEntered(ch, es, true)
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
				es := target.FindEntrustedShopByOwner(ch.GetID())
				ch.GameWorld.GetDispatchSystem().CallCharacter(ch.GetID(), func(ctx actor.Context, c *Character) {
					if c == nil {
						return
					}
					if es == nil {
						c.Listener.OnEntrustedShopCheck(c, pconst.EntrustedShopRemoteLocation, StoreBankNothingMapID, StoreBankNothingChannel)
						return
					}
					c.remoteShop = remoteShop{sn: es.OID, mapID: target.TemplateID()}
					c.Listener.OnEntrustedShopRemoteVisit(c, es.OID)
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
		return es.Visit(ch)
	}
	if ps, ok := m.GetObject(constant.ObjectTypePersonalShop, sn).(*PersonalShop); ok {
		return ps.Visit(ch)
	}
	if ch.remoteShop.sn != sn {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}

	remote := ch.remoteShop
	ch.remoteShop = remoteShop{}
	target := ch.GameWorld.GetMapSystem().Get(remote.mapID)
	if target == nil {
		return &MiniRoomEnterError{Code: pconst.MiniRoomEnterClosed}
	}
	ch.GameWorld.GetMapSystem().Call(target, func(actor.Context) {
		es, ok := target.GetObject(constant.ObjectTypeEntrustedShop, sn).(*EntrustedShop)
		ch.GameWorld.GetDispatchSystem().CallCharacter(ch.GetID(), func(ctx actor.Context, c *Character) {
			if c == nil {
				return
			}
			if ok == false || es.OwnerID != c.GetID() {
				c.Listener.OnMiniRoomEnterFailed(c, pconst.MiniRoomEnterClosed)
				return
			}
			c.RejectMiniRoom(es.Visit(c))
		})
	})
	return nil
}
